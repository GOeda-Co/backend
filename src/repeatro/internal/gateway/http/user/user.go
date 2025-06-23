package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"repeatro/src/repeatro/internal/gateway"
	// "repeatro/src/user/pkg/scheme"
)

// Gateway layer replaces repository layer in the services where other microservices cannot be used
// Inside specific methods i do requests to my other microservices
type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr}
}

func (g *Gateway) Login(ctx context.Context, body io.ReadCloser) (string, error) {
	req, err := http.NewRequest(http.MethodPost, g.addr+"/login", body)
	if err != nil {
		return "", err
	}

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
		// fmt.Println("HERE")
		return "", fmt.Errorf("%s", response.Status)
	}

	var token string
	if err := json.NewDecoder(response.Body).Decode(&token); err != nil {
		return "", err
	}
	return token, nil
}

func (g *Gateway) Register(ctx context.Context, body io.ReadCloser) (string, error) {
	req, err := http.NewRequest(http.MethodPost, g.addr+"/register", body)
	if err != nil {
		return "", err
	}

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

	var token string
	if err := json.NewDecoder(response.Body).Decode(&token); err != nil {
		return "", err
	}
	return token, nil
}
