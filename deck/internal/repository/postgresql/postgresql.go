package postgresql

import (
	"log/slog"

	models "github.com/tomatoCoderq/deck/pkg/model"

	"github.com/google/uuid"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	
)

type Repository struct {
	db *gorm.DB
}

func New(connectionString string, log *slog.Logger) *Repository {
	db, err := gorm.Open(postgres.Open(connectionString))
	if err != nil {
		log.Error("Error during opening database")
	}

	db.AutoMigrate(&models.Deck{})

	return &Repository{db: db}
}

func (r *Repository) AddDeck(deck *models.Deck) error {
	return r.db.Create(deck).Error
}

func (r *Repository) ReadAllDecksOfUser(userId uuid.UUID) ([]models.Deck, error) {
	var decks []models.Deck
	err := r.db.Where("created_by = ?", userId).Preload("Cards").Find(&decks).Error
	return decks, err
}

func (r *Repository) ReadAllDecks() ([]models.Deck, error) {
	var decks []models.Deck
	err := r.db.Preload("Cards").Find(&decks).Error
	return decks, err
}

func (r *Repository) ReadDeck(deckId uuid.UUID) (*models.Deck, error) {
	var deck models.Deck
	err := r.db.Where("deck_id = ?", deckId).Preload("Cards").First(&deck).Error
	if err != nil {
		return nil, err
	}
	return &deck, nil
}

func (r *Repository) DeleteDeck(deckId uuid.UUID) error {
	return r.db.Delete(&models.Deck{}, "deck_id = ?", deckId).Error
}

func (r *Repository) FindAllCardsInDeck(deckId uuid.UUID) ([]models.Card, error) {
	var cards []models.Card
	err := r.db.Where("deck_id = ?", deckId).Find(&cards).Error
	return cards, err
}

func (r *Repository) AddCardToDeck(cardId uuid.UUID, deckId uuid.UUID) error {
	return r.db.Model(&models.Card{}).Where("card_id = ?", cardId).Update("deck_id", deckId).Error
}
