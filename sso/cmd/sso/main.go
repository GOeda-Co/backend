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
			"sso/config/config.yaml",
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

	// Initialize logger
	log := initLogger(cfg.Env)
	log.Info("Initializing logger")
	log.Debug("Initializing debug mode")

	// Initialize app
	fmt.Println(cfg.TokenTTL)
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
