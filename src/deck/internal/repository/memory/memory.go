package memory

import (
	"errors"
	"sync"

	cardModel "repeatro/src/card/pkg/model"
	"repeatro/src/deck/pkg/model"

	"github.com/google/uuid"
)

type Repository struct {
	mu    sync.RWMutex
	decks map[uuid.UUID]*models.Deck
	cards map[uuid.UUID]*cardModel.Card
}

func NewInMemoryRepo() *Repository {
	return &Repository{
		decks: make(map[uuid.UUID]*models.Deck),
		cards: make(map[uuid.UUID]*cardModel.Card),
	}
}

func (r *Repository) AddDeck(deck *models.Deck) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if deck.DeckId == uuid.Nil {
		deck.DeckId = uuid.New()
	}
	r.decks[deck.DeckId] = deck
	return nil
}

func (r *Repository) ReadAllDecksOfUser(userId uuid.UUID) ([]models.Deck, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []models.Deck
	for _, d := range r.decks {
		if d.CreatedBy == userId {
			result = append(result, *d)
		}
	}
	return result, nil
}

func (r *Repository) ReadAllDecks() ([]models.Deck, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []models.Deck
	for _, d := range r.decks {
		result = append(result, *d)
	}
	return result, nil
}

func (r *Repository) ReadDeck(deckId uuid.UUID) (*models.Deck, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	deck, exists := r.decks[deckId]
	if !exists {
		return nil, errors.New("deck not found")
	}
	return deck, nil
}

func (r *Repository) DeleteDeck(deckId uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.decks[deckId]; !exists {
		return errors.New("deck not found")
	}
	delete(r.decks, deckId)
	return nil
}

func (r *Repository) FindAllCardsInDeck(deckId uuid.UUID) ([]cardModel.Card, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []cardModel.Card
	for _, c := range r.cards {
		if c.DeckID == deckId {
			result = append(result, *c)
		}
	}
	return result, nil
}

func (r *Repository) AddCardToDeck(cardId uuid.UUID, deckId uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	card, exists := r.cards[cardId]
	if !exists {
		return errors.New("card not found")
	}
	card.DeckID = deckId
	return nil
}
