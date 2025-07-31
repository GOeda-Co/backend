package app

import (
	"log/slog"

	"github.com/tomatoCoderq/deck/internal/app/grpc"
	// client "github.com/tomatoCoderq/deck/internal/clients/sso/grpc"
	"github.com/tomatoCoderq/deck/internal/lib/security"
	"github.com/tomatoCoderq/deck/internal/repository/postgresql"
	"github.com/tomatoCoderq/deck/internal/services/deck"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storageAddress string,
	security security.Security,
) *App {
	storage := postgresql.New(storageAddress, log)

	authService := services.New(log, storage)
	grpcApp := grpcapp.New(log, authService, grpcPort, security)

	return &App{
		GRPCServer: grpcApp,
	}
}
