package postgresql

import (
	// "fmt"
	"log/slog"

	// model "github.com/tomatoCoderq/deck/pkg/model"
	modelCard "github.com/GOeda-Co/proto-contract/model/card"
	model "github.com/GOeda-Co/proto-contract/model/deck"
	"github.com/tomatoCoderq/deck/migrations"

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

	migrations.MigrateToLatest(db, log)

	if err = db.AutoMigrate(&model.Deck{}); err != nil {
		log.Error("Error during auto migration", "error", err)
	}

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

func (r *Repository) SearchAllPublicDecks() ([]model.Deck, error) {
	var decks []model.Deck
	err := r.db.Where("is_public = ?", true).Preload("Cards").Find(&decks).Error
	return decks, err
}

func (r *Repository) SearchUserPublicDecks(userId uuid.UUID) ([]model.Deck, error) {
	var decks []model.Deck
	err := r.db.Where("is_public = ? AND created_by = ?", true, userId).Preload("Cards").Find(&decks).Error
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
	var deck *model.Deck
	var err error

	tx := r.db.Begin()

	if deck, err = r.ReadDeck(deckId); err != nil {
		tx.Rollback()
		return err
	}

	// After adding a card to a deck, deck's is_public updated to deck's is_public
	if err := tx.Model(&modelCard.Card{}).
		Where("card_id = ?", cardId).
		Update("deck_id", deckId).
		Update("is_public", deck.IsPublic).Error; err != nil {
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
