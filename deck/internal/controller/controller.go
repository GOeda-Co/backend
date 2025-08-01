package controller

import (
	"github.com/google/uuid"
	// models "github.com/tomatoCoderq/deck/pkg/model"
	modelCard "github.com/GOeda-Co/proto-contract/model/card"
	"github.com/GOeda-Co/proto-contract/model/deck"
)

type Deck interface {
	AddDeck(deck *model.Deck) (*model.Deck, error)
	ReadAllDecksOfUser(userId uuid.UUID) ([]model.Deck, error)
	ReadAllCardsFromDeck(deckId uuid.UUID, userId uuid.UUID) ([]modelCard.Card, error)
	SearchAllPublicDecks() ([]model.Deck, error)
	SearchUserPublicDecks(userId string) ([]model.Deck, error)
	ReadDeck(deckId uuid.UUID, userId uuid.UUID) (*model.Deck, error)
	DeleteDeck(deckId uuid.UUID, userId uuid.UUID) error
	AddCardToDeck(cardId uuid.UUID, deckId uuid.UUID, userId uuid.UUID) error
}
