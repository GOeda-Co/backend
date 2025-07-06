package app

import (
	"log/slog"

	"github.com/tomatoCoderq/card/internal/app/grpc"
	client "github.com/tomatoCoderq/card/internal/clients/sso/grpc"
	"github.com/tomatoCoderq/card/internal/lib/security"
	"github.com/tomatoCoderq/card/internal/repository/postgresql"
	"github.com/tomatoCoderq/card/internal/services/card"
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

	authService := services.New(log, storage)
	grpcApp := grpcapp.New(log, authService, grpcPort, client, security)

	return &App {
		GRPCServer: grpcApp,
	}
}