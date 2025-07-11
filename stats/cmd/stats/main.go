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

	client "github.com/tomatoCoderq/stats/internal/clients/sso/grpc"
	"github.com/tomatoCoderq/stats/internal/config"
	"github.com/tomatoCoderq/stats/internal/lib/security"

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
	cfg := config.MustLoad()

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
