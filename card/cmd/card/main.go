package main

import (
	// "context"
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"log/slog"
	"os"

	"github.com/joho/godotenv"
	// ssoClient "github.com/tomatoCoderq/card/internal/clients/sso/grpc"
	statClient "github.com/tomatoCoderq/card/internal/clients/stats/grpc"
	"github.com/tomatoCoderq/card/internal/config"
	"github.com/tomatoCoderq/card/internal/lib/security"
	"gopkg.in/yaml.v3"

	app "github.com/tomatoCoderq/card/internal/app"
	// "github.com/gin-gonic/gin"
)

// const serviceName = "cards"

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	_ = godotenv.Load(".env")

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		possiblePaths := []string{
			"config/config.yaml",
			"../../config/config.yaml",
			"/app/config/config.yaml",
			"card/config/config.yaml",
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				break
			}
		}
		if configPath == "" {
			panic(fmt.Errorf("could not find config file in any of the expected locations"))
		}
	}

	// Read and expand config file
	raw, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Errorf("could not read config file %s: %w", configPath, err))
	}
	expanded := os.ExpandEnv(string(raw))

	var cfg config.Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		panic(fmt.Errorf("could not unmarshal config file %s: %w", configPath, err))
	}

	log := setupLogger(cfg.Env)

	log.Info(
		"starting...",
		slog.String("env", cfg.Env),
	)
	log.Debug("debug messages are enabled")

	statClient, err := statClient.New(context.Background(), log, cfg.Clients.STAT.Address, cfg.Clients.STAT.Timeout.Abs(), cfg.Clients.STAT.RetriesCount)
	if err != nil {
		panic(err)
	}

	security := security.Security{
		PrivateKey:      cfg.Secret,
		ExpirationDelta: 600 * time.Minute,
	}

	application := app.New(log, cfg.GRPC.Port, cfg.ConnectionString, statClient, security) // ssoClient, statClient commented out
	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	log.Info("Gracefully stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, nil))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
