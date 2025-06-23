package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"net/http"
	_ "net/http/pprof"

	// "repeatro/internal/middlewares"
	"repeatro/internal/security"
	"repeatro/src/pkg"
	"repeatro/src/pkg/discovery/consul"
	httphandler "repeatro/src/repeatro/internal/controller/http"
	"repeatro/src/repeatro/internal/gateway/http/card"
	"repeatro/src/repeatro/internal/gateway/http/deck"
	"repeatro/src/repeatro/internal/gateway/http/user"
	"repeatro/src/repeatro/internal/service"

	"github.com/gin-gonic/gin"
)

const serviceName = "cards"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
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

	log.Println("Starting the repeatro service")
	cardGateway := card.New(registry)
	userGateway := user.New(registry)
	deckGateway := deck.New(registry)
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

	
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}
