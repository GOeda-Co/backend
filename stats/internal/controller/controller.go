package controller

import (
	statsv1 "github.com/GOeda-Co/proto-contract/gen/go/stats"
)

type Service interface {
	GetAverageGrade(uid, deckId string, timeRange statsv1.TimeRange) (float64, error)
	GetCardsReviewedCount(uid, deckId string, timeRange statsv1.TimeRange) (int32, error)
	// GetCardsLearnedCount(uid, deckId string, timeRange statsv1.TimeRange) (int32, error)
}