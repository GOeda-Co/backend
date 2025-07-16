package model

import (
	"time"

	// "repeatro/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Deck struct {
	DeckId        uuid.UUID `gorm:"type:uuid;primaryKey;" json:"deck_id"`
	CreatedBy     uuid.UUID `gorm:"references:UserId;constraint:OnDelete:CASCADE;" json:"created_by"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	Name          string    `gorm:"type:varchar(100);not null;default:null" json:"name"`
	Description   string    `gorm:"type:varchar(100);" json:"description"`
	CardsQuantity uint      `gorm:"type:int unsigned;default=0" json:"cards_quantity"`
	Cards         []Card    `gorm:"foreignKey:CardId;constraint:OnDelete:CASCADE"`
}

func (d *Deck) BeforeCreate(tx *gorm.DB) error {
	d.DeckId = uuid.New()
	return nil
}
