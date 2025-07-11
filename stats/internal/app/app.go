package app

import (
	"log/slog"

	"github.com/tomatoCoderq/stats/internal/app/grpc"
	client "github.com/tomatoCoderq/stats/internal/clients/sso/grpc"
	"github.com/tomatoCoderq/stats/internal/lib/security"
	"github.com/tomatoCoderq/stats/internal/repository/postgresql"
	"github.com/tomatoCoderq/stats/internal/service/stats"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
    log *slog.Logger,
    grpcPort int,
    storageAddress string,
	client *client.Client,
	security security.Security,
) *App {
	storage := postgresql.New(storageAddress, log)

	statsService := stats.New(log, storage)
	grpcApp := grpcapp.New(log, statsService, grpcPort, client, security)

	return &App {
		GRPCServer: grpcApp,
	}
}