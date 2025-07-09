package app

import (
	"log/slog"

	"github.com/tomatoCoderq/repeatro/internal/lib/security"
	httpApp "github.com/tomatoCoderq/repeatro/internal/app/http"
	ssoClient "github.com/tomatoCoderq/repeatro/internal/clients/sso/grpc"
	cardClient "github.com/tomatoCoderq/repeatro/internal/clients/card/grpc"

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
	security security.Security,
) *App {
	grpcApp := httpApp.New(log, port, address, ssoClient, cardClient, security)

	return &App {
		HttpServer: grpcApp,
	}
}