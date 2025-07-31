package controller

import (
	"time"

	statsv1 "github.com/GOeda-Co/proto-contract/gen/go/stats"
	"github.com/google/uuid"
)

type Service interface {
	GetAverageGrade(uid, deckId string, timeRange statsv1.TimeRange) (float64, error)
	GetCardsReviewedCount(uid, deckId string, timeRange statsv1.TimeRange) (int32, error)
	AddRecord(uid uuid.UUID, deckId, dardId string, CreatedAt time.Time, grade int) (string, error)
	// GetCardsLearnedCount(uid, deckId string, timeRange statsv1.TimeRange) (int32, error)
}
