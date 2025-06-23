package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"net/http"
	_ "net/http/pprof"
	"repeatro/internal/config"
	"repeatro/internal/security"
	deckHttp "repeatro/src/deck/internal/controller/http"
	"repeatro/src/deck/internal/repository/postgresql"
	"repeatro/src/deck/internal/service"
	"repeatro/src/pkg"
	"repeatro/src/pkg/discovery/consul"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

const serviceName = "decks"

func main() {
	var port int
	flag.IntVar(&port, "port", 8085, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	security := security.Security{ExpirationDelta: 600 * time.Minute}
	security.GetKyes()

	config := config.InitConfig("config")

	log.Println("Starting the user service")
	repo := postgresql.NewPostgresRepo(config, newLogger)
	service := services.CreateNewService(repo)
	ctrl := deckHttp.CreateNewController(service)

	default_router := gin.Default()

	router := default_router.Group("/decks")

	router.Handle(http.MethodPost, "", ctrl.AddDeck)
	router.Handle(http.MethodGet, "", ctrl.ReadAllDecks)
	router.Handle(http.MethodGet, "/:id", ctrl.ReadDeck)
	router.Handle(http.MethodDelete, "/:id", ctrl.DeleteDeck)
	router.Handle(http.MethodPost, "/:deck_id/cards/:card_id", ctrl.AddCardToDeck) // post one card
	router.Handle(http.MethodGet, "/:id/cards", ctrl.ReadCardsFromDeck)

	if err := default_router.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}
