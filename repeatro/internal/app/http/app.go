package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	ssoClient "github.com/tomatoCoderq/repeatro/internal/clients/sso/grpc"
	cardClient "github.com/tomatoCoderq/repeatro/internal/clients/card/grpc"
	deckClient "github.com/tomatoCoderq/repeatro/internal/clients/deck/grpc"
	httpRepeatro "github.com/tomatoCoderq/repeatro/internal/controller/http"
	"github.com/tomatoCoderq/repeatro/internal/lib/security"
)

type App struct {
	log        *slog.Logger
	port       int
	httpServer *http.Server
}

func New(
	log *slog.Logger,
	port int,
	address string,
	ssoClient *ssoClient.Client,
	cardClient *cardClient.Client,
	deckClient *deckClient.Client,
	security security.Security,
) *App {
	router := gin.Default()
	router.Use(gin.Recovery())

	ctrl := httpRepeatro.New(ssoClient, cardClient, deckClient)
	router.Handle(http.MethodPost, "/register", ctrl.Register)
	router.Handle(http.MethodPost, "/login", ctrl.Login)

	//TODO CHANGE 

	cards := router.Group("/cards")
	cards.Use(security.AuthMiddleware())

	cards.Handle(http.MethodPost, "", ctrl.AddCard)
	cards.Handle(http.MethodGet, "", ctrl.ReadAllCardsToLearn)
	cards.Handle(http.MethodPut, "/:id", ctrl.UpdateCard)
	cards.Handle(http.MethodDelete, "/:id", ctrl.DeleteCard)
	cards.Handle(http.MethodPost, "/answers", ctrl.AddAnswers)

	decks := router.Group("/decks")
	decks.Use(security.AuthMiddleware())
	decks.Handle(http.MethodPost, "", ctrl.AddDeck)
	decks.Handle(http.MethodGet, "", ctrl.ReadAllDecks)
	decks.Handle(http.MethodGet, "/:id", ctrl.ReadDeck)
	decks.Handle(http.MethodDelete, "/:id", ctrl.DeleteDeck)
	// decks.Handle(http.MethodPut, "/:id", ctrl.UpdateCard)
	decks.Handle(http.MethodPost, "/:deck_id/cards/:card_id", ctrl.AddCardToDeck)
	decks.Handle(http.MethodGet, "/:id/cards", ctrl.ReadCardsFromDeck)
	

	httpServer := &http.Server{
		Addr:    address,
		Handler: router,
	}

	return &App{
		log:        log,
		port:       port,
		httpServer: httpServer,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "app.Run"

	a.log.Info("HTTP server started", slog.String("addr", a.httpServer.Addr))

	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: http server error: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "app.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping HTTP server", slog.Int("port", a.port))

	if err := a.httpServer.Shutdown(context.Background()); err != nil {
		a.log.Error("HTTP server shutdown error", slog.String("op", op), slog.String("error", err.Error()))
	}
}