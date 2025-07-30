package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	// "fmt"
	"log/slog"
	"os"

	// "time"

	// "net/http"

	"github.com/joho/godotenv"
	client "github.com/tomatoCoderq/stats/internal/clients/sso/grpc"
	"github.com/tomatoCoderq/stats/internal/config"
	"github.com/tomatoCoderq/stats/internal/lib/security"
	"gopkg.in/yaml.v3"

	// userHttp "github.com/tomatoCoderq/stats/internal/controller/http"
	// "github.com/tomatoCoderq/stats/internal/lib/security"
	// "github.com/tomatoCoderq/stats/internal/repository/postgresql"
	// services "github.com/tomatoCoderq/stats/internal/services/stats"

	app "github.com/tomatoCoderq/stats/internal/app"
	// "github.com/gin-gonic/gin"
)

const serviceName = "stats"

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Load .env variables from multiple possible locations
	_ = godotenv.Load("/app/.env")
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../../.env")

	// Determine config file path - check environment variable first, then fallback to defaults
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// Try multiple possible locations for config file
		possiblePaths := []string{
			"config/config.yaml",
			"../../config/config.yaml",
			"/app/config/config.yaml",
			"stats/config/config.yaml",
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

	ssoClient, err := client.New(context.Background(), log, cfg.Clients.SSO.Address, cfg.Clients.SSO.Timeout.Abs(), cfg.Clients.SSO.RetriesCount)
	if err != nil {
		panic(err)
	}

	fmt.Println(ssoClient)

	security := security.Security{
		PrivateKey:      cfg.Secret,
		ExpirationDelta: 600 * time.Minute,
	}

	application := app.New(log, cfg.GRPC.Port, cfg.ConnectionString, ssoClient, security)
	go func() {
		application.GRPCServer.MustRun()
	}()

	//TODO: Завершить работу программы
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	//TODO: Add close for db
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
