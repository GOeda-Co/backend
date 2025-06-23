package card

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"repeatro/src/card/pkg/model"
	"repeatro/src/repeatro/internal/gateway"
	// "github.com/google/uuid"
	// "github.com/golang-jwt/jwt/v5"
	// "github.com/gin-gonic/gin"
	// "github.com/google/uuid"
)

// Gateway layer replaces repository layer in the services where other microservices cannot be used
// Inside specific methods i do requests to my other microservices
type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr}
}

func (g *Gateway) GetCards(ctx context.Context, userClaims string) ([]*model.Card, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/cards", nil)
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

	fmt.Println(g.addr + "/cards")
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
	req, err := http.NewRequest(http.MethodPost, g.addr+"/cards", body)
	if err != nil {
		return nil, err
	}

	fmt.Println(g.addr + "/cards")

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
	req, err := http.NewRequest(http.MethodPut, g.addr+"/cards/"+cardId, body)
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

	fmt.Println(g.addr + "/cards")
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
	req, err := http.NewRequest(http.MethodDelete, g.addr+"/cards/"+cardId, nil)
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

	fmt.Println(g.addr + "/cards")
	if response.StatusCode == http.StatusNotFound {
		return "", gateway.ErrNotFound
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s", response.Status)
	}

	return response.Status, nil
}

func (g *Gateway) AddAnswers(ctx context.Context, userClaims string, body io.ReadCloser) (string, error) {
	req, err := http.NewRequest(http.MethodPost, g.addr+"/cards/answers", body)
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

	fmt.Println(g.addr + "/cards/answers")
	if response.StatusCode == http.StatusNotFound {
		return "", gateway.ErrNotFound
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s", response.Status)
	}

	return response.Status, nil
}
