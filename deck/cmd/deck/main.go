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
	client "github.com/tomatoCoderq/deck/internal/clients/sso/grpc"
	"github.com/tomatoCoderq/deck/internal/config"
	"github.com/tomatoCoderq/deck/internal/lib/security"
	"gopkg.in/yaml.v3"

	// userHttp "github.com/tomatoCoderq/card/internal/controller/http"
	// "github.com/tomatoCoderq/card/internal/lib/security"
	// "github.com/tomatoCoderq/card/internal/repository/postgresql"
	// services "github.com/tomatoCoderq/card/internal/services/card"

	app "github.com/tomatoCoderq/deck/internal/app"
	// "github.com/gin-gonic/gin"
)

const serviceName = "cards"

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

	// TODO: Завершить работу программы
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	// TODO: Add close for db
	log.Info("Gracefully stopped")

	// storage := postgresql.New(cfg.ConnectionString, log)
	// service := services.New(log, storage)
	// ctrl := userHttp.CreateNewController(service)

	// security := security.Security{
	// 	PrivateKey:      cfg.Secret,
	// 	ExpirationDelta: 600 * time.Minute,
	// }

	// default_router := gin.Default()
	// default_router.Use(security.AuthMiddleware())
	// default_router.Use(security.IsAdminMiddleware(*ssoClient))

	// router := default_router.Group("/cards")

	// fmt.Println("EHEbqhe")

	// router.Handle(http.MethodPost, "", ctrl.AddCard)
	// router.Handle(http.MethodGet, "", ctrl.ReadAllCardsToLearn)
	// router.Handle(http.MethodPut, "/:id", ctrl.UpdateCard)
	// router.Handle(http.MethodDelete, "/:id", ctrl.DeleteCard)
	// router.Handle(http.MethodPost, "/answers", ctrl.AddAnswers)

	// srv := &http.Server{
	// 	Addr:         cfg.Address,
	// 	Handler:      default_router,
	// 	ReadTimeout:  cfg.HTTPServer.Timeout,
	// 	WriteTimeout: cfg.HTTPServer.Timeout,
	// 	IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	// }

	// log.Info("Starting server")

	// if err := srv.ListenAndServe(); err != nil {
	// 	panic(err)
	// }

}

// start app

// end app

// logger

// TODO: technically each microservice should have separated main and current one should be divided into three
// func main() {
// 	var port int
// 	flag.IntVar(&port, "port", 8084, "API handler port")
// 	flag.Parse()
// 	log.Printf("Starting the movie service on port %d", port)
// 	registry, err := consul.NewRegistry("localhost:8500")
// 	if err != nil {
// 		panic(err)
// 	}
// 	ctx := context.Background()
// 	instanceID := discovery.GenerateInstanceID(serviceName)
// 	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
// 		panic(err)
// 	}
// 	go func() {
// 		for {
// 			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
// 				log.Println("Failed to report healthy state: " + err.Error())
// 			}
// 			time.Sleep(1 * time.Second)
// 		}
// 	}()
// 	defer registry.Deregister(ctx, instanceID, serviceName)

// 	newLogger := logger.New(
// 		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
// 		logger.Config{
// 			SlowThreshold:             time.Second,   // Slow SQL threshold
// 			LogLevel:                  logger.Silent, // Log level
// 			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
// 			ParameterizedQueries:      true,          // Don't include params in the SQL log
// 			Colorful:                  false,         // Disable color
// 		},
// 	)
// 	security := security.Security{ExpirationDelta: 600 * time.Minute}
// 	security.GetKyes()

// 	config := config.InitConfig("config")

// 	log.Println("Starting the user service")
// 	repo := postgresql.NewPostgresRepo(config, newLogger)
// 	service := services.CreateNewService(repo)
// 	ctrl := userHttp.CreateNewController(service)

// 	default_router := gin.Default()
// 	// default_router.Use(security.AuthMiddleware())

// 	router := default_router.Group("/cards")

// 	router.Handle(http.MethodPost, "", ctrl.AddCard)
// 	router.Handle(http.MethodGet, "", ctrl.ReadAllCardsToLearn)
// 	router.Handle(http.MethodPut, "/:id", ctrl.UpdateCard)
// 	router.Handle(http.MethodDelete, "/:id", ctrl.DeleteCard)
// 	router.Handle(http.MethodPost, "/answers", ctrl.AddAnswers)

// 	if err := default_router.Run(fmt.Sprintf(":%d", port)); err != nil {
// 		panic(err)
// 	}
// }

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
