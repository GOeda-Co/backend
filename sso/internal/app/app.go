package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/storage/postgresql"
	"sso/internal/services/auth"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}


func New(
    log *slog.Logger,
    grpcPort int,
    storageAddress string,
    tokenTTL time.Duration,
) *App {
	storage, err := postgresql.New(storageAddress, log)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App {
		GRPCServer: grpcApp,
	}
}