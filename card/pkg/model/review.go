package model

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ResultId  uuid.UUID
	UserID    uuid.UUID
	CardID    uuid.UUID
	DeckId    uuid.UUID
	CreatedAt time.Time
	Grade     int32
}
