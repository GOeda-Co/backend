package memory

import (
	"errors"
	"sync"
	"time"

	"repeatro/src/card/pkg/model"
	"repeatro/src/card/pkg/scheme"

	"github.com/google/uuid"
)

func updateCardFields(cardInitial *model.Card, card *schemes.UpdateCardScheme) {
	if card.Word != "" {
		cardInitial.Word = card.Word
	}
	if card.Translation != "" {
		cardInitial.Translation = card.Translation
	}
	if card.Easiness != 0 {
		cardInitial.Easiness = card.Easiness
	}
	if !card.UpdatedAt.IsZero() {
		cardInitial.UpdatedAt = card.UpdatedAt
	}
	if card.Interval != 0 {
		cardInitial.Interval = card.Interval
	}
	if !card.ExpiresAt.IsZero() {
		cardInitial.ExpiresAt = card.ExpiresAt
	}
	if card.RepetitionNumber != 0 {
		cardInitial.RepetitionNumber = card.RepetitionNumber
	}
	// TODO: Add tags here
}

type Repository struct {
	mu    sync.RWMutex
	cards map[uuid.UUID]*model.Card
}

func NewInMemoryCardRepo() *Repository {
	return &Repository{
		cards: make(map[uuid.UUID]*model.Card),
	}
}

func (r *Repository) AddCard(card *model.Card) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if card.CardId == uuid.Nil {
		card.CardId = uuid.New()
	}
	r.cards[card.CardId] = card
	return nil
}

func (r *Repository) ReadAllCards(userId uuid.UUID) ([]model.Card, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.Card
	for _, c := range r.cards {
		if c.CreatedBy == userId && c.ExpiresAt.Before(time.Now()) {
			result = append(result, *c)
		}
	}
	return result, nil
}

func (r *Repository) ReadCard(cardId uuid.UUID) (*model.Card, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	card, exists := r.cards[cardId]
	if !exists {
		return nil, errors.New("card not found")
	}
	return card, nil
}

func (r *Repository) UpdateCard(card *model.Card, cardUpdate *schemes.UpdateCardScheme) (*model.Card, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.cards[card.CardId]; !exists {
		return nil, errors.New("card not found")
	}
	updateCardFields(card, cardUpdate)
	return card, nil
}

func (r *Repository) PureUpdate(card *model.Card) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.cards[card.CardId]; !exists {
		return errors.New("card not found")
	}
	r.cards[card.CardId] = card
	return nil
}

func (r *Repository) DeleteCard(cardId uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.cards[cardId]; !exists {
		return errors.New("card not found")
	}
	delete(r.cards, cardId)
	return nil
}
