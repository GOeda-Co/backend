package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ResultId  uuid.UUID `gorm:"type:uuid;primaryKey;"`
	UserID    uuid.UUID
	CardID    uuid.UUID
	DeckId    uuid.UUID
	CreatedAt time.Time
	Grade     int32
}

func (Review) TableName() string {
	return "results"
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ResultId == uuid.Nil {
		r.ResultId = uuid.New()
	}
	return nil
}
