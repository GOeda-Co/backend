package main

import (
	"log"
	"time"

	"net/http"
	_ "net/http/pprof"

	// "repeatro/internal/middlewares"
	"repeatro/internal/security"
	httphandler "repeatro/src/repeatro/internal/controller/http"
	"repeatro/src/repeatro/internal/gateway/http/card"
	"repeatro/src/repeatro/internal/gateway/http/deck"
	"repeatro/src/repeatro/internal/gateway/http/user"
	"repeatro/src/repeatro/internal/service"

	"github.com/gin-gonic/gin"
)

// // TODO: technically each microservice should have separated main and current one should be divided into three
// func main() {
// 	config := config.InitConfig("config")

// 	security := security.Security{ExpirationDelta: 600 * time.Minute}
// 	security.GetKyes()

// 	db := repositories.InitDatabase(config)

// 	repositories.InitGooseMigration(db)

// 	server := server.InitHTTPServer(config, db, security)

// 	go func() {
// 		log.Println(http.ListenAndServe("localhost:6060", nil))
// 	}()

// 	server.StartHttpServer()

// }

func main() {
	log.Println("Starting the repeatro service")
	cardGateway := card.New("http://localhost:8084")
	userGateway := user.New("http://localhost:8082")
	deckGateway := deck.New("http://localhost:8085")
	service := service.New(cardGateway, userGateway, deckGateway)
	ctrl := httphandler.New(service)

	security := security.Security{ExpirationDelta: 600 * time.Minute}
	security.GetKyes()

	router := gin.Default()
	router.Handle(http.MethodPost, "/login", ctrl.Login)
	router.Handle(http.MethodPost, "/register", ctrl.Register)
	
	secured := router.Group("")
	secured.Use(security.AuthMiddleware())

	cards := secured.Group("/cards")
	decks := secured.Group("/decks")
	
	//TODO: Too many repetitions of handlers (need to write tool func with func passing as arg)

	cards.Handle(http.MethodGet, "", ctrl.GetAllCards)
	cards.Handle(http.MethodPost, "", ctrl.AddCard)
	cards.Handle(http.MethodPut, "/:id", ctrl.UpdateCard)
	cards.Handle(http.MethodDelete, "/:id", ctrl.DeleteCard)
	cards.Handle(http.MethodPost, "/answers", ctrl.AddAnswers)


	decks.Handle(http.MethodPost, "", ctrl.AddDeck)
	decks.Handle(http.MethodGet, "", ctrl.ReadAllDecks)
	decks.Handle(http.MethodGet, "/:id", ctrl.ReadDeck)
	decks.Handle(http.MethodDelete, "/:id", ctrl.DeleteDeck)
	decks.Handle(http.MethodPost, "/:deck_id/cards/:card_id", ctrl.AddCardToDeck) // post one card
	decks.Handle(http.MethodGet, "/:id/cards", ctrl.ReadCardsFromDeck)

	
	if err := router.Run(":8083"); err != nil {
		panic(err)
	}
}
