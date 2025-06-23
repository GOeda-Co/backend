package main

import (
	"log"
	"os"
	"time"

	_ "net/http/pprof"
	"repeatro/internal/config"
	"repeatro/internal/security"
	"net/http"
	deckHttp "repeatro/src/deck/internal/controller/http"
	"repeatro/src/deck/internal/repository/postgresql"
	"repeatro/src/deck/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

// TODO: technically each microservice should have separated main and current one should be divided into three
func main() {
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


	if err := default_router.Run(":8085"); err != nil {
		panic(err)
	}
}
