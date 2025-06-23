package services

import (
	"errors"
	"repeatro/src/deck/pkg/model"
	"repeatro/src/card/pkg/model"
	// "repeatro/src/deck/internal/repository/postgresql"

	"github.com/google/uuid"
)

var ErrUnauthorized = errors.New("you do not own this deck")

type deckRepository interface {
	AddDeck(deck *models.Deck) error
	ReadAllDecksOfUser(userId uuid.UUID) ([]models.Deck, error)
	ReadAllDecks() ([]models.Deck, error)
	ReadDeck(deckId uuid.UUID) (*models.Deck, error)
	DeleteDeck(deckId uuid.UUID) error
	AddCardToDeck(cardId uuid.UUID, deckId uuid.UUID) error
	FindAllCardsInDeck(deckId uuid.UUID) ([]model.Card, error )
}

type Service struct {
	deckRepository deckRepository
}

func CreateNewService(deckRepository deckRepository) *Service {
	return &Service{
		deckRepository: deckRepository,
	}
}

// type DeckServiceInterface interface {
// 	AddCard(deck *models.Deck, userId uuid.UUID) (*models.Deck, error)
// 	ReadAllDecksOfUser(userId uuid.UUID) ([]models.Deck, error)
// 	ReadAllCardsFromDeck(deckId uuid.UUID, userId uuid.UUID) ([]models.Card, error)
// 	ReadDeck(deckId uuid.UUID, userId uuid.UUID) (*models.Deck, error)
// 	DeleteDeck(deckId uuid.UUID, userId uuid.UUID) error
// 	AddCardToDeck(cardId uuid.UUID, deckId uuid.UUID, userId uuid.UUID) error
// }

func (ds *Service) AddCard(deck *models.Deck, userId uuid.UUID) (*models.Deck, error) {
	deck.CreatedBy = userId
	err := ds.deckRepository.AddDeck(deck)
	if err != nil {
		return nil, err
	}
	return deck, nil
}

func (ds *Service) ReadAllDecksOfUser(userId uuid.UUID) ([]models.Deck, error) {
	return ds.deckRepository.ReadAllDecksOfUser(userId)
}

func (ds *Service) ReadDeck(deckId uuid.UUID, userId uuid.UUID) (*models.Deck, error) {
	deck, err := ds.deckRepository.ReadDeck(deckId)
	if err != nil {
		return nil, err
	}

	if deck.CreatedBy != userId {
		return nil, ErrUnauthorized
	}
	return deck, nil
}

func (ds *Service) ReadAllCardsFromDeck(deckId uuid.UUID, userId uuid.UUID) ([]model.Card, error) {
	cards, err := ds.deckRepository.FindAllCardsInDeck(deckId)
	if err != nil {
		return nil, err
	}
	return cards, nil
	
}
  
func (ds *Service) DeleteDeck(deckId uuid.UUID, userId uuid.UUID) error {
	deck, err := ds.deckRepository.ReadDeck(deckId)
	if err != nil {
		return err
	}
	if deck.CreatedBy != userId {
		return ErrUnauthorized
	}
	return ds.deckRepository.DeleteDeck(deckId)
}

func (ds *Service) AddCardToDeck(cardId uuid.UUID, deckId uuid.UUID, userId uuid.UUID) error {
	deck, err := ds.deckRepository.ReadDeck(deckId)
	if err != nil {
		return err
	}
	if deck.CreatedBy != userId {
		return ErrUnauthorized
	}
	return ds.deckRepository.AddCardToDeck(cardId, deckId)
}
