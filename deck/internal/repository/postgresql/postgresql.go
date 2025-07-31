package postgresql

import (
	// "fmt"
	"log/slog"

	// model "github.com/tomatoCoderq/deck/pkg/model"
	model "github.com/GOeda-Co/proto-contract/model/deck"
	modelCard "github.com/GOeda-Co/proto-contract/model/card"

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

	db.AutoMigrate(&model.Deck{})

	return &Repository{db: db}
}

func (r *Repository) AddDeck(deck *model.Deck) error {
	return r.db.Create(deck).Error
}

func (r *Repository) ReadAllDecksOfUser(userId uuid.UUID) ([]model.Deck, error) {
	var decks []model.Deck
	err := r.db.Where("created_by = ?", userId).Preload("Cards").Find(&decks).Error
	return decks, err
}

func (r *Repository) ReadAllDecks() ([]model.Deck, error) {
	var decks []model.Deck
	err := r.db.Preload("Cards").Find(&decks).Error
	return decks, err
}

func (r *Repository) ReadDeck(deckId uuid.UUID) (*model.Deck, error) {
	var deck model.Deck
	err := r.db.Where("deck_id = ?", deckId).Preload("Cards").First(&deck).Error
	if err != nil {
		return nil, err
	}
	return &deck, nil
}

func (r *Repository) DeleteDeck(deckId uuid.UUID) error {
	return r.db.Delete(&model.Deck{}, "deck_id = ?", deckId).Error
}

func (r *Repository) FindAllCardsInDeck(deckId uuid.UUID) ([]modelCard.Card, error) {
	var cards []modelCard.Card
	err := r.db.Where("deck_id = ?", deckId).Find(&cards).Error
	return cards, err
}

func (r *Repository) AddCardToDeck(cardId uuid.UUID, deckId uuid.UUID) error {
	tx := r.db.Begin()

	if err := tx.Model(&modelCard.Card{}).
		Where("card_id = ?", cardId).
		Update("deck_id", deckId).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.Deck{}).
		Where("deck_id = ?", deckId).
		UpdateColumn("cards_quantity", gorm.Expr("cards_quantity + ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
	// return r.db.Model(&model.Card{}).Where("card_id = ?", cardId).Update("deck_id", deckId).Error
}
