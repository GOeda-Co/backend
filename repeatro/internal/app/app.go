package app

import (
	"log/slog"

	httpApp "github.com/tomatoCoderq/repeatro/internal/app/http"
	cardClient "github.com/tomatoCoderq/repeatro/internal/clients/card/grpc"
	deckClient "github.com/tomatoCoderq/repeatro/internal/clients/deck/grpc"
	ssoClient "github.com/tomatoCoderq/repeatro/internal/clients/sso/grpc"
	statClient "github.com/tomatoCoderq/repeatro/internal/clients/stats/grpc"
	"github.com/tomatoCoderq/repeatro/internal/lib/security"
)

type App struct {
	HttpServer *httpApp.App
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
	grpcApp := httpApp.New(log, port, address, ssoClient, cardClient, deckClient, statClient, security)

	return &App{
		HttpServer: grpcApp,
	}
}
