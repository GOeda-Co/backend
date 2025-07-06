package services

import (
	"errors"
	
	models "github.com/tomatoCoderq/deck/pkg/model"
	
	// "repeatro/src/deck/internal/repository/postgresql"

	"github.com/google/uuid"
)

var ErrUnauthorized = errors.New("you do not own this deck")

type DeckRepository interface {
	AddDeck(deck *models.Deck) error
	ReadAllDecksOfUser(userId uuid.UUID) ([]models.Deck, error)
	ReadAllDecks() ([]models.Deck, error)
	ReadDeck(deckId uuid.UUID) (*models.Deck, error)
	DeleteDeck(deckId uuid.UUID) error
	AddCardToDeck(cardId uuid.UUID, deckId uuid.UUID) error
	FindAllCardsInDeck(deckId uuid.UUID) ([]models.Card, error )
}

type Service struct {
	DeckRepository DeckRepository
}

func CreateNewService(DeckRepository DeckRepository) *Service {
	return &Service{
		DeckRepository: DeckRepository,
	}
}

type DeckServiceInterface interface {
	AddCard(deck *models.Deck, userId uuid.UUID) (*models.Deck, error)
	ReadAllDecksOfUser(userId uuid.UUID) ([]models.Deck, error)
	ReadAllCardsFromDeck(deckId uuid.UUID, userId uuid.UUID) ([]models.Card, error)
	ReadDeck(deckId uuid.UUID, userId uuid.UUID) (*models.Deck, error)
	DeleteDeck(deckId uuid.UUID, userId uuid.UUID) error
	AddCardToDeck(cardId uuid.UUID, deckId uuid.UUID, userId uuid.UUID) error
}

func (ds *Service) AddCard(deck *models.Deck, userId uuid.UUID) (*models.Deck, error) {
	deck.CreatedBy = userId
	err := ds.DeckRepository.AddDeck(deck)
	if err != nil {
		return nil, err
	}
	return deck, nil
}

func (ds *Service) ReadAllDecksOfUser(userId uuid.UUID) ([]models.Deck, error) {
	return ds.DeckRepository.ReadAllDecksOfUser(userId)
}

func (ds *Service) ReadDeck(deckId uuid.UUID, userId uuid.UUID) (*models.Deck, error) {
	deck, err := ds.DeckRepository.ReadDeck(deckId)
	if err != nil {
		return nil, err
	}

	if deck.CreatedBy != userId {
		return nil, ErrUnauthorized
	}
	return deck, nil
}

func (ds *Service) ReadAllCardsFromDeck(deckId uuid.UUID, userId uuid.UUID) ([]models.Card, error) {
	cards, err := ds.DeckRepository.FindAllCardsInDeck(deckId)
	if err != nil {
		return nil, err
	}
	return cards, nil
	
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
