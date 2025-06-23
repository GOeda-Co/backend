package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"repeatro/src/card/pkg/model"
	deckModel "repeatro/src/deck/pkg/model"
	"repeatro/src/repeatro/internal/gateway"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	// schemes "repeatro/src/user/pkg/scheme"
	// "github.com/google/uuid"
)

//here i define

var ErrDeckNotFound = errors.New("deck not found")

//card, deck, user interfaces for gateaways

//and service struct

func GetUserClaimsBytes(ctx *gin.Context) ([]byte, error) {
	userClaims, ok := ctx.Get("userClaims")
	if !ok {
		return nil, fmt.Errorf("some error")
	}
	userClaimsValid, ok := userClaims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("some error2")
	}

	userClaimsBytes, err := json.Marshal(userClaimsValid)
	if err != nil {
		return nil, err
	}
	return userClaimsBytes, nil
}

var ErrNotFound = errors.New("card not found")

type cardGateway interface {
	GetCards(ctx context.Context, userClaims string) ([]*model.Card, error)
	AddCard(ctx context.Context, userClaims string, body io.ReadCloser) (*model.Card, error)
	UpdateCard(ctx context.Context, userClaims string, body io.ReadCloser, cardId string) (*model.Card, error)
	DeleteCard(ctx context.Context, userClaims string, cardId string) (string, error)
	AddAnswers(ctx context.Context, userClaims string, body io.ReadCloser) (string, error)
}

type deckGateway interface {
	AddDeck(ctx context.Context, userClaims string, body io.ReadCloser) (*deckModel.Deck, error)
	ReadAllDecks(ctx context.Context, userClaims string) ([]*deckModel.Deck, error)
	ReadDeck(ctx context.Context, userClaims string, deckId string) (*deckModel.Deck, error)
	ReadCardsFromDeck(ctx context.Context, userClaims string, deckId string) ([]*model.Card, error)
	DeleteDeck(ctx context.Context, userClaims string, deckId string) error
	AddCardToDeck(ctx context.Context, userClaims string, deckId string, cardId string) error
}

type userGateway interface {
	Login(ctx context.Context, body io.ReadCloser) (string, error)
	Register(ctx context.Context, body io.ReadCloser) (string, error)
}
type Service struct {
	cardGateway cardGateway
	userGateway userGateway
	deckGateway deckGateway
}

func New(cardGateway cardGateway, userGateway userGateway, deckGateway deckGateway) *Service {
	return &Service{cardGateway, userGateway, deckGateway}
}

func (s *Service) GetCards(ctx context.Context, ctxGin *gin.Context) ([]*model.Card, error) {
	userClaims, ok := ctxGin.Get("userClaims")
	if !ok {
		return nil, fmt.Errorf("some error")
	}
	userClaimsValid, ok := userClaims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("some error2")
	}

	userClaimsBytes, err := json.Marshal(userClaimsValid)
	if err != nil {
		return nil, err
	}

	cards, err := s.cardGateway.GetCards(ctx, string(userClaimsBytes))
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return cards, nil
}

func (s *Service) AddCard(ctx context.Context, ctxGin *gin.Context) (*model.Card, error) {
	userClaimsBytes, err := GetUserClaimsBytes(ctxGin)
	if err != nil {
		return nil, err
	}

	card, err := s.cardGateway.AddCard(ctx, string(userClaimsBytes), ctxGin.Request.Body)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return card, nil
}

func (s *Service) UpdateCard(ctx context.Context, ctxGin *gin.Context) (*model.Card, error) {
	userClaimsBytes, err := GetUserClaimsBytes(ctxGin)
	if err != nil {
		return nil, err
	}

	cardId := ctxGin.Param("id")
	if cardId == "" {
		return nil, fmt.Errorf("no id provided")
	}

	card, err := s.cardGateway.UpdateCard(ctx, string(userClaimsBytes), ctxGin.Request.Body, cardId)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return card, nil
}

func (s *Service) DeleteCard(ctx context.Context, ctxGin *gin.Context) (string, error) {
	userClaimsBytes, err := GetUserClaimsBytes(ctxGin)
	if err != nil {
		return "", err
	}

	cardId := ctxGin.Param("id")
	if cardId == "" {
		return "", fmt.Errorf("no id provided")
	}

	statusCode, err := s.cardGateway.DeleteCard(ctx, string(userClaimsBytes), cardId)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return "", ErrNotFound
	} else if err != nil {
		return "", err
	}
	return statusCode, nil
}

func (s *Service) AddAnswers(ctx context.Context, ctxGin *gin.Context) (string, error) {
	userClaimsBytes, err := GetUserClaimsBytes(ctxGin)
	if err != nil {
		return "", err
	}

	statusCode, err := s.cardGateway.AddAnswers(ctx, string(userClaimsBytes), ctxGin.Request.Body)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return "", ErrNotFound
	} else if err != nil {
		return "", err
	}
	return statusCode, nil
}

func (s *Service) AddDeck(ctx context.Context, ctxGin *gin.Context) (*deckModel.Deck, error) {
	userClaimsBytes, err := GetUserClaimsBytes(ctxGin)
	if err != nil {
		return nil, err
	}

	deck, err := s.deckGateway.AddDeck(ctx, string(userClaimsBytes), ctxGin.Request.Body)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrDeckNotFound
	}
	return deck, err
}

func (s *Service) ReadAllDecksOfUser(ctx context.Context, ctxGin *gin.Context) ([]*deckModel.Deck, error) {
	userClaimsBytes, err := GetUserClaimsBytes(ctxGin)
	if err != nil {
		return nil, err
	}

	decks, err := s.deckGateway.ReadAllDecks(ctx, string(userClaimsBytes))
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrDeckNotFound
	}
	return decks, err
}

func (s *Service) ReadDeck(ctx context.Context, ctxGin *gin.Context, deckId string) (*deckModel.Deck, error) {
	userClaimsBytes, err := GetUserClaimsBytes(ctxGin)
	if err != nil {
		return nil, err
	}

	deck, err := s.deckGateway.ReadDeck(ctx, string(userClaimsBytes), deckId)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrDeckNotFound
	}
	return deck, err
}

func (s *Service) ReadAllCardsFromDeck(ctx context.Context, ctxGin *gin.Context, deckId string) ([]*model.Card, error) {
	userClaimsBytes, err := GetUserClaimsBytes(ctxGin)
	if err != nil {
		return nil, err
	}

	cards, err := s.deckGateway.ReadCardsFromDeck(ctx, string(userClaimsBytes), deckId)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrDeckNotFound
	}
	return cards, err
}

func (s *Service) DeleteDeck(ctx context.Context, ctxGin *gin.Context, deckId string) error {
	userClaimsBytes, err := GetUserClaimsBytes(ctxGin)
	if err != nil {
		return err
	}

	err = s.deckGateway.DeleteDeck(ctx, string(userClaimsBytes), deckId)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return ErrDeckNotFound
	}
	return err
}

func (s *Service) AddCardToDeck(ctx context.Context, ctxGin *gin.Context, deckId string, cardId string) error {
	userClaimsBytes, err := GetUserClaimsBytes(ctxGin)
	if err != nil {
		return err
	}

	err = s.deckGateway.AddCardToDeck(ctx, string(userClaimsBytes), deckId, cardId)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return ErrDeckNotFound
	}
	return err
}

func (s *Service) Login(ctx context.Context, ctxGin *gin.Context) (string, error) {
	token, err := s.userGateway.Login(ctx, ctxGin.Request.Body)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Register(ctx context.Context, ctxGin *gin.Context) (string, error) {
	token, err := s.userGateway.Register(ctx, ctxGin.Request.Body)
	if err != nil {
		return "", err
	}

	return token, nil
}
