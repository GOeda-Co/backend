package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/tomatoCoderq/repeatro/docs"
	cardClient "github.com/tomatoCoderq/repeatro/internal/clients/card/grpc"
	deckClient "github.com/tomatoCoderq/repeatro/internal/clients/deck/grpc"
	ssoClient "github.com/tomatoCoderq/repeatro/internal/clients/sso/grpc"
	statClient "github.com/tomatoCoderq/repeatro/internal/clients/stats/grpc"
	"github.com/tomatoCoderq/repeatro/internal/config"
	"github.com/tomatoCoderq/repeatro/internal/lib/security"
	"gopkg.in/yaml.v3"

	app "github.com/tomatoCoderq/repeatro/internal/app"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

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
			"repeatro/config/config.yaml",
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

	ssoClient, err := ssoClient.New(context.Background(), log, cfg.Clients.SSO.Address, cfg.Clients.SSO.Timeout.Abs(), cfg.Clients.SSO.RetriesCount)
	if err != nil {
		panic(err)
	}

	cardClient, err := cardClient.New(context.Background(), log, cfg.Clients.CARD.Address, cfg.Clients.CARD.Timeout.Abs(), cfg.Clients.CARD.RetriesCount)
	if err != nil {
		panic(err)
	}

	deckClient, err := deckClient.New(context.Background(), log, cfg.Clients.DECK.Address, cfg.Clients.DECK.Timeout.Abs(), cfg.Clients.DECK.RetriesCount)
	if err != nil {
		panic(err)
	}

	statClient, err := statClient.New(context.Background(), log, cfg.Clients.STAT.Address, cfg.Clients.SSO.Timeout.Abs(), cfg.Clients.STAT.RetriesCount)
	if err != nil {
		panic(err)
	}

	security := security.Security{
		PrivateKey:      cfg.Secret,
		ExpirationDelta: 600 * time.Minute,
	}

	application := app.New(log, cfg.HTTPServer.Port, cfg.HTTPServer.Address, ssoClient, cardClient, deckClient, statClient, security)
	go func() {
		application.HttpServer.MustRun()
	}()

	// Creating a channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Wait for a signal to stop the application
	<-stop

	application.HttpServer.Stop()
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
