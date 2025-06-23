package deck

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	model "repeatro/src/deck/pkg/model"
	modelCard "repeatro/src/card/pkg/model"
	"repeatro/src/repeatro/internal/gateway"
)

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr}
}

func (g *Gateway) AddDeck(ctx context.Context, userClaims string, body io.ReadCloser) (*model.Deck, error) {
	req, err := http.NewRequest(http.MethodPost, g.addr+"/decks", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("userClaims", userClaims)
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	var deck *model.Deck
	if err := json.NewDecoder(resp.Body).Decode(&deck); err != nil {
		return nil, err
	}

	return deck, nil
}

func (g *Gateway) ReadAllDecks(ctx context.Context, userClaims string) ([]*model.Deck, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/decks", nil)
	if err != nil {
		return nil, err
	}
	
	fmt.Println(g.addr+"/decks")

	req.Header.Set("userClaims", userClaims)
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	var decks []*model.Deck
	if err := json.NewDecoder(resp.Body).Decode(&decks); err != nil {
		return nil, err
	}

	return decks, nil
}

func (g *Gateway) ReadDeck(ctx context.Context, userClaims string, deckId string) (*model.Deck, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/decks/"+deckId, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("userClaims", userClaims)
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	var deck *model.Deck
	if err := json.NewDecoder(resp.Body).Decode(&deck); err != nil {
		return nil, err
	}

	return deck, nil
}

func (g *Gateway) ReadCardsFromDeck(ctx context.Context, userClaims string, deckId string) ([]*modelCard.Card, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/decks/"+deckId+"/cards", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("userClaims", userClaims)
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	var cards []*modelCard.Card
	if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
		return nil, err
	}

	return cards, nil
}

func (g *Gateway) DeleteDeck(ctx context.Context, userClaims string, deckId string) error {
	req, err := http.NewRequest(http.MethodDelete, g.addr+"/decks/"+deckId, nil)
	if err != nil {
		return err
	}

	req.Header.Set("userClaims", userClaims)
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return gateway.ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}

	return nil
}

func (g *Gateway) AddCardToDeck(ctx context.Context, userClaims string, deckId string, cardId string) error {
	req, err := http.NewRequest(http.MethodPost, g.addr+"/decks/"+deckId+"/cards/"+cardId, nil)
	if err != nil {
		return err
	}

	req.Header.Set("userClaims", userClaims)
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return gateway.ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}

	return nil
}
