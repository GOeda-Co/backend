//	@title			Repeatro
//	@version		1.0
//	@description	Repeatro Swagger describes all endpoints.

//	@host		localhost:8080
//	@BasePath	/

// @contact.name	khasan
// @contact.email	xasanFN@mail.ru
// @contact.url	https://t.me/tomatocoder

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	cardClient "github.com/tomatoCoderq/repeatro/internal/clients/card/grpc"
	deckClient "github.com/tomatoCoderq/repeatro/internal/clients/deck/grpc"
	ssoClient "github.com/tomatoCoderq/repeatro/internal/clients/sso/grpc"
	statClient "github.com/tomatoCoderq/repeatro/internal/clients/stats/grpc"
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
	statClient *statClient.Client,
	security security.Security,
) *App {
	router := gin.Default()
	router.Use(gin.Recovery(), cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // adjust for your frontend
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	ctrl := httpRepeatro.New(log, ssoClient, cardClient, deckClient, statClient)
	router.Handle(http.MethodPost, "/register", ctrl.Register)
	router.Handle(http.MethodPost, "/login", ctrl.Login)
	router.Handle(http.MethodGet, "/admin", ctrl.IsAdmin)

	//TODO CHANGE

	cards := router.Group("/cards")
	cards.Use(security.AuthMiddleware())

	cards.Handle(http.MethodPost, "", ctrl.AddCard)
	cards.Handle(http.MethodGet, "/learn", ctrl.ReadAllCardsToLearn)
	cards.Handle(http.MethodGet, "", ctrl.ReadAllCards)
	cards.Handle(http.MethodGet, "/search", ctrl.SearchPublicCards)
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

	stats := router.Group("/stats")
	stats.Use(security.AuthMiddleware())

	stats.Handle(http.MethodGet, "/average", ctrl.GetAverageGrade)
	stats.Handle(http.MethodGet, "/count", ctrl.GetCardsReviewedCount)

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
