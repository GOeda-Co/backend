package card

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"repeatro/src/card/pkg/model"
	"repeatro/src/pkg"
	"repeatro/src/repeatro/internal/gateway"
	// "github.com/google/uuid"
	// "github.com/golang-jwt/jwt/v5"
	// "github.com/gin-gonic/gin"
	// "github.com/google/uuid"
)



// Gateway layer replaces repository layer in the services where other microservices cannot be used
// Inside specific methods i do requests to my other microservices
type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

func (g *Gateway) GetCards(ctx context.Context, userClaims string) ([]*model.Card, error) {
	url, err := gateway.GetAvailableAddresses(ctx, "cards", "/cards", g.registry.ServiceAddresses)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("userClaims", userClaims)

	req = req.WithContext(ctx)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", response.Status)
	}

	var cards []*model.Card
	if err := json.NewDecoder(response.Body).Decode(&cards); err != nil {
		return nil, err
	}
	return cards, nil
}

func (g *Gateway) AddCard(ctx context.Context, userClaims string, body io.ReadCloser) (*model.Card, error) {
	url, err := gateway.GetAvailableAddresses(ctx, "cards", "/cards", g.registry.ServiceAddresses)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("userClaims", userClaims)
	req = req.WithContext(ctx)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", response.Status)
	}

	var card *model.Card
	if err := json.NewDecoder(response.Body).Decode(&card); err != nil {
		return nil, err
	}
	return card, nil
}

func (g *Gateway) UpdateCard(ctx context.Context, userClaims string, body io.ReadCloser, cardId string) (*model.Card, error) {
	url, err := gateway.GetAvailableAddresses(ctx, "cards", "/cards", g.registry.ServiceAddresses)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("userClaims", userClaims)

	req = req.WithContext(ctx)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", response.Status)
	}

	var card *model.Card
	if err := json.NewDecoder(response.Body).Decode(&card); err != nil {
		return nil, err
	}
	return card, nil
}

func (g *Gateway) DeleteCard(ctx context.Context, userClaims string, cardId string) (string, error) {
	url, err := gateway.GetAvailableAddresses(ctx, "cards", "/cards/"+cardId, g.registry.ServiceAddresses)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("userClaims", userClaims)

	req = req.WithContext(ctx)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return "", gateway.ErrNotFound
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s", response.Status)
	}

	return response.Status, nil
}

func (g *Gateway) AddAnswers(ctx context.Context, userClaims string, body io.ReadCloser) (string, error) {
	url, err := gateway.GetAvailableAddresses(ctx, "cards", "/cards/answers", g.registry.ServiceAddresses)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("userClaims", userClaims)

	req = req.WithContext(ctx)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return "", gateway.ErrNotFound
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s", response.Status)
	}

	return response.Status, nil
}
