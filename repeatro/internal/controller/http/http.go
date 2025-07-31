package http

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/swaggo/swag/example/celler/httputil"
	cardClient "github.com/tomatoCoderq/repeatro/internal/clients/card/grpc"
	deckClient "github.com/tomatoCoderq/repeatro/internal/clients/deck/grpc"
	ssoClient "github.com/tomatoCoderq/repeatro/internal/clients/sso/grpc"
	statClient "github.com/tomatoCoderq/repeatro/internal/clients/stats/grpc"
)

type Controller struct {
	log        *slog.Logger
	ssoClient  *ssoClient.Client
	cardClient *cardClient.Client
	deckClient *deckClient.Client
	statClient *statClient.Client
}

func New(log *slog.Logger, ssoClient *ssoClient.Client, cardClient *cardClient.Client, deckClient *deckClient.Client, statClient *statClient.Client) *Controller {
	return &Controller{
		log:        log,
		ssoClient:  ssoClient,
		cardClient: cardClient,
		deckClient: deckClient,
		statClient: statClient,
	}
}

func GetUserIdFromContext(ctx *gin.Context) (uuid.UUID, error) {
	userClaims, exists := ctx.Get("userClaims")
	if !exists {
		return uuid.UUID{}, fmt.Errorf("user claims do not exist")
	}

	claimsMap, ok := userClaims.(jwt.MapClaims)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("cannot convert claims to map")
	}

	fmt.Println("VIUST")

	userIdString, ok := claimsMap["uid"].(string)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("cannot get user_id from claims map")
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("cannot parse uuid")
	}

	return userId, nil
}

func GetIsAdminFromContext(ctx *gin.Context) (bool, error) {
	userClaims, exists := ctx.Get("userClaims")
	if !exists {
		return false, fmt.Errorf("user claims do not exist")
	}

	claimsMap, ok := userClaims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("cannot convert claims to map")
	}

	slog.Debug("Available claims", "claims", claimsMap)

	// Check if admin field exists
	if adminVal, exists := claimsMap["admin"]; exists {
		if isAdmin, ok := adminVal.(bool); ok {
			return isAdmin, nil
		} else {
			slog.Debug("Admin field exists but is not boolean", "value", adminVal)
			return false, fmt.Errorf("admin field is not boolean type")
		}
	} else {
		slog.Debug("Admin field does not exist in claims, defaulting to false")
		return false, nil
	}
}

// type ErrorResponse struct {
// 	Error string `json:"error"`
// }
