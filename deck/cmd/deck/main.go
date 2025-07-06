package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tomatoCoderq/deck/internal/config"
	userHttp "github.com/tomatoCoderq/deck/internal/controller/http"
	"github.com/tomatoCoderq/deck/internal/lib/security"
	"github.com/tomatoCoderq/deck/internal/repository/postgresql"
	"github.com/tomatoCoderq/deck/internal/service"
)

const serviceName = "decks"

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
		slog.String("version", "123"),
	)
	log.Debug("debug messages are enabled")

	storage := postgresql.New(cfg.ConnectionString, log)
	service := services.CreateNewService(storage)
	ctrl := userHttp.CreateNewController(service)

	security := security.Security{ExpirationDelta: 600 * time.Minute}
	security.GetKyes(cfg.Secret)

	default_router := gin.Default()
	default_router.Use(security.AuthMiddleware())


	router := default_router.Group("/decks")

	router.Handle(http.MethodPost, "", ctrl.AddDeck)
	router.Handle(http.MethodGet, "", ctrl.ReadAllDecks)
	router.Handle(http.MethodGet, "/:id", ctrl.ReadDeck)
	router.Handle(http.MethodDelete, "/:id", ctrl.DeleteDeck)
	router.Handle(http.MethodPost, "/:deck_id/cards/:card_id", ctrl.AddCardToDeck) // post one card
	router.Handle(http.MethodGet, "/:id/cards", ctrl.ReadCardsFromDeck)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      default_router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("Starting server")

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

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

// func main() {
// 	var port int
// 	flag.IntVar(&port, "port", 8085, "API handler port")
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
// 	ctrl := deckHttp.CreateNewController(service)

// 	default_router := gin.Default()

// 	router := default_router.Group("/decks")

// 	router.Handle(http.MethodPost, "", ctrl.AddDeck)
// 	router.Handle(http.MethodGet, "", ctrl.ReadAllDecks)
// 	router.Handle(http.MethodGet, "/:id", ctrl.ReadDeck)
// 	router.Handle(http.MethodDelete, "/:id", ctrl.DeleteDeck)
// 	router.Handle(http.MethodPost, "/:deck_id/cards/:card_id", ctrl.AddCardToDeck) // post one card
// 	router.Handle(http.MethodGet, "/:id/cards", ctrl.ReadCardsFromDeck)

// 	if err := default_router.Run(fmt.Sprintf(":%d", port)); err != nil {
// 		panic(err)
// 	}
// }
