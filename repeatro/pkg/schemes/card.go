package schemes

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AnswerScheme struct {
	CardId uuid.UUID `json:"card_id"`
	Grade  int       `json:"grade"`
}

type UpdateCardScheme struct {
	Word             string         `json:"word"`
	Translation      string         `json:"translation"`
	Easiness         float64        `json:"easiness"`
	UpdatedAt        time.Time      `json:"updated_at"`
	Interval         int            `json:"interval"`
	ExpiresAt        time.Time      `json:"expires_at"`
	RepetitionNumber int            `json:"repetition_number"`
	Tags             pq.StringArray `json:"tags"`
}
