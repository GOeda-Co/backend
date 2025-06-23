package postgresql

import (
	"log"
	"time"

	"repeatro/src/card/pkg/model"
	"repeatro/src/card/pkg/scheme"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

func NewPostgresRepo(config *viper.Viper, newLogger logger.Interface) *Repository {
	db, err := gorm.Open(postgres.Open(config.GetString("database.connection_string")), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatalf("Error during opening database")
	}

	db.AutoMigrate(&model.Card{})

	return &Repository{db: db}
}

type CardRepositoryInterface interface {
	AddCard(card *model.Card) error
	ReadAllCards(userId uuid.UUID) ([]model.Card, error)
	ReadCard(cardId uuid.UUID) (*model.Card, error)
	PureUpdate(card *model.Card) error
	UpdateCard(card *model.Card, cardUpdate *schemes.UpdateCardScheme) (*model.Card, error)
	DeleteCard(cardId uuid.UUID) error
}

func (cr Repository) AddCard(card *model.Card) error {
	return cr.db.Create(card).Error
}

func (cr Repository) ReadAllCards(userId uuid.UUID) ([]model.Card, error) {
	var cards []model.Card
	err := cr.db.Where("expires_at < ?", time.Now()).Where("created_by = ?", userId).Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, err
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
