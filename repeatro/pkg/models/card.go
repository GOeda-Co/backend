package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// NOTE: I don't use rn gorm.Model due to redundancy of some fields

type Card struct {
	CardId           uuid.UUID      `gorm:"type:uuid;primaryKey;" json:"card_id"`
	CreatedBy        uuid.UUID      `gorm:"references:UserId;constraint:OnDelete:CASCADE;" json:"created_by"`
	User             User           `gorm:"foreignKey:CreatedBy;references:UserId" json:"-"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	Word             string         `gorm:"type:varchar(100);not null;default:null" json:"word"`
	Translation      string         `gorm:"type:varchar(100);not null;default:null" json:"translation"`
	Easiness         float64        `gorm:"type:float64;not null;default:2.5" json:"easiness"`
	UpdatedAt        time.Time      `gorm:"type:autoCreateTime" json:"updated_at"`
	Interval         int            `gorm:"type:smallint;default=0" json:"interval"`
	ExpiresAt        time.Time      `gorm:"autoCreateTime" json:"expires_at"`
	RepetitionNumber int            `gorm:"type:smallint;default=0" json:"repetition_number"`
	DeckID           uuid.UUID      `gorm:"type:uuid;index" json:"deck_id"`
	Tags             pq.StringArray `gorm:"type:text[]" json:"tags"`
	IsPublic         bool           `gorm:"default:false" json:"is_public"`
}

func (c *Card) BeforeCreate(tx *gorm.DB) error {
	c.CardId = uuid.New()
	c.ExpiresAt = time.Now().Add(10 * time.Second)
	return nil
}
