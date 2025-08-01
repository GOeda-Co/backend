package services

import (
	"errors"
	"log/slog"

	
	"github.com/GOeda-Co/proto-contract/model/deck"
	modelCard "github.com/GOeda-Co/proto-contract/model/card"

	// "repeatro/src/deck/internal/repository/postgresql"

	"github.com/google/uuid"
)

var ErrUnauthorized = errors.New("you do not own this deck")

type DeckRepository interface {
	AddDeck(deck *model.Deck) error
	ReadAllDecksOfUser(userId uuid.UUID) ([]model.Deck, error)
	ReadAllDecks() ([]model.Deck, error)
	ReadDeck(deckId uuid.UUID) (*model.Deck, error)
	SearchAllPublicDecks() ([]model.Deck, error)
	SearchUserPublicDecks(userId uuid.UUID) ([]model.Deck, error)
	DeleteDeck(deckId uuid.UUID) error
	AddCardToDeck(cardId uuid.UUID, deckId uuid.UUID) error
	FindAllCardsInDeck(deckId uuid.UUID) ([]modelCard.Card, error)
}

type Service struct {
	DeckRepository DeckRepository
}

func New(log *slog.Logger, DeckRepository DeckRepository) *Service {
	return &Service{
		DeckRepository: DeckRepository,
	}
}

func (ds *Service) AddDeck(deck *model.Deck) (*model.Deck, error) {
	err := ds.DeckRepository.AddDeck(deck)
	if err != nil {
		return nil, err
	}
	return deck, nil
}

func (ds *Service) ReadAllDecksOfUser(userId uuid.UUID) ([]model.Deck, error) {
	return ds.DeckRepository.ReadAllDecksOfUser(userId)
}

func (ds *Service) ReadDeck(deckId uuid.UUID, userId uuid.UUID) (*model.Deck, error) {
	deck, err := ds.DeckRepository.ReadDeck(deckId)
	if err != nil {
		return nil, err
	}

	if deck.CreatedBy != userId {
		return nil, ErrUnauthorized
	}
	return deck, nil
}

func (ds *Service) ReadAllCardsFromDeck(deckId uuid.UUID, userId uuid.UUID) ([]modelCard.Card, error) {
	cards, err := ds.DeckRepository.FindAllCardsInDeck(deckId)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (ds *Service) SearchAllPublicDecks() ([]model.Deck, error) {
	return ds.DeckRepository.SearchAllPublicDecks()
}

func (ds *Service) SearchUserPublicDecks(userId string) ([]model.Deck, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	return ds.DeckRepository.SearchUserPublicDecks(userUUID)
}

func (ds *Service) DeleteDeck(deckId uuid.UUID, userId uuid.UUID) error {
	deck, err := ds.DeckRepository.ReadDeck(deckId)
	if err != nil {
		return err
	}
	if deck.CreatedBy != userId {
		return ErrUnauthorized
	}
	return ds.DeckRepository.DeleteDeck(deckId)
}

func (ds *Service) AddCardToDeck(cardId uuid.UUID, deckId uuid.UUID, userId uuid.UUID) error {
	deck, err := ds.DeckRepository.ReadDeck(deckId)
	if err != nil {
		return err
	}
	if deck.CreatedBy != userId {
		return ErrUnauthorized
	}
	return ds.DeckRepository.AddCardToDeck(cardId, deckId)
}
