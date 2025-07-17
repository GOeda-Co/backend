package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/config"
	"sso/internal/app"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: инициализировать объект конфига
	config := config.MustLoad()

	// TODO: инициализировать логгер
	log := initLogger(config.Env)
	log.Info("Initializing logger")
	log.Debug("Initializing debug mode")

	// TODO: инициализировать приложение (app)
	application := app.New(log, config.GRPC.Address, config.ConnectionString, config.TokenTTL)

	// TODO: запустить gRPC-сервер приложения
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

func initLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
