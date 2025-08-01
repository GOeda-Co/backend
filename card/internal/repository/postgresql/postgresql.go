package postgresql

import (
	"log/slog"
	"time"

	// "github.com/tomatoCoderq/card/pkg/model"
	"github.com/GOeda-Co/proto-contract/model/card"
	// "github.com/tomatoCoderq/card/pkg/scheme"
	schemes "github.com/GOeda-Co/proto-contract/scheme/card"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	db *gorm.DB
}

func New(connectionString string, log *slog.Logger) *Repository {
	db, err := gorm.Open(postgres.Open(connectionString))
	if err != nil {
		log.Error("Error during opening database")
		return nil
	}

	if err := db.AutoMigrate(&model.Card{}); err != nil {
		log.Error("Error during auto migration", "error", err)
		return nil
	}

	return &Repository{db: db}
}

func (cr Repository) AddCard(card *model.Card) error {
	return cr.db.Create(card).Error
}

func (cr Repository) ReadAllOwnCardsToLearn(userId uuid.UUID) ([]model.Card, error) {
	var cards []model.Card
	err := cr.db.
		Where("expires_at < ?", time.Now()).
		Where("created_by = ?", userId).
		Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, err
}

func (cr Repository) ReadAllOwnCards(userId uuid.UUID) ([]model.Card, error) {
	var cards []model.Card
	err := cr.db.
		Where("created_by = ?", userId).
		Find(&cards).
		Error
	if err != nil {
		return nil, err
	}
	return cards, err
}

func (cr Repository) SearchAllPublicCards() ([]model.Card, error) {
	var cards []model.Card
	err := cr.db.
		Where("is_public = ?", true).
		Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (cr Repository) SearchUserPublicCards(userId uuid.UUID) ([]model.Card, error) {
	var cards []model.Card
	err := cr.db.
		Where("is_public = ?", true).
		Where("created_by = ?", userId).
		Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (cr Repository) ReadCard(cardId uuid.UUID) (*model.Card, error) {
	var card model.Card
	err := cr.db.Where("card_id = ?", cardId).Find(&card).Error
	return &card, err
}

func (cr Repository) UpdateCard(card *model.Card, cardUpdate *schemes.UpdateCardScheme) (*model.Card, error) {
	updateCardFields(card, cardUpdate)
	return card, cr.db.Updates(card).Error
}

func (cr Repository) PureUpdate(card *model.Card) error {
	return cr.db.Updates(card).Error
}

func (cr Repository) DeleteCard(cardId uuid.UUID) error {
	err := cr.db.Delete(&model.Card{}, "card_id = ?", cardId).Error
	return err
}
