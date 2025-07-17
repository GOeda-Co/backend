package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"sso/config"
	"sso/internal/app"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Load .env variables
	_ = godotenv.Load("/app/.env")

	// Read and expand config.yaml
	raw, err := os.ReadFile("/app/config/config.yaml")
	if err != nil {
		panic(fmt.Errorf("could not read config.yaml: %w", err))
	}
	expanded := os.ExpandEnv(string(raw))

	var cfg config.Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		panic(fmt.Errorf("could not unmarshal config.yaml: %w", err))
	}

	// Initialize logger
	log := initLogger(cfg.Env)
	log.Info("Initializing logger")
	log.Debug("Initializing debug mode")

	// Initialize app
	application := app.New(log, cfg.GRPC.Port, cfg.ConnectionString, cfg.TokenTTL)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	application.GRPCServer.Stop()
	log.Info("Gracefully stopped")
}

func initLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
