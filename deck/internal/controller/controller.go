package controller

import (
	"github.com/google/uuid"
	models "github.com/tomatoCoderq/deck/pkg/model"
)


type Deck interface {
	AddDeck(deck *models.Deck) (*models.Deck, error)
	ReadAllDecksOfUser(userId uuid.UUID) ([]models.Deck, error)
	ReadAllCardsFromDeck(deckId uuid.UUID, userId uuid.UUID) ([]models.Card, error)
	ReadDeck(deckId uuid.UUID, userId uuid.UUID) (*models.Deck, error)
	DeleteDeck(deckId uuid.UUID, userId uuid.UUID) error
	AddCardToDeck(cardId uuid.UUID, deckId uuid.UUID, userId uuid.UUID) error
}