package stats

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	statsv1 "github.com/GOeda-Co/proto-contract/gen/go/stats"
	"github.com/google/uuid"
)

func calculateTimeRange(tr statsv1.TimeRange) (start, end time.Time) {
	now := time.Now()
	switch tr {
	case statsv1.TimeRange_DAILY:
		start = now.AddDate(0, 0, -1)
	case statsv1.TimeRange_WEEKLY:
		start = now.AddDate(0, 0, -7)
	case statsv1.TimeRange_MONTHLY:
		start = now.AddDate(0, -1, 0)
	default:
		start = now
	}
	end = now

	return
}

type Repository interface {
	AverageGrade(uid, deckId uuid.UUID, startTime, endTime time.Time) (float64, error)
	CountReviewedCards(uid, deckId uuid.UUID, startTime, endTime time.Time) (int32, error)
	AddRecord(uid, deckId, cardId uuid.UUID, createdAt time.Time, grade int) (string, error)
	// GetCardsLearnedCount(uid, cardId string, startTime, endTime time.Time) (int32, error)
}

type Service struct {
	log  *slog.Logger
	repo Repository
}

func New(log *slog.Logger, repo Repository) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}

func (s *Service) GetAverageGrade(uid, deckId string, timeRange statsv1.TimeRange) (float64, error) {
	if timeRange == statsv1.TimeRange_TIME_RANGE_UNSPECIFIED {
		return 0, errors.New("time range must be specified")
	}

	fmt.Println("POINT3")

	startTime, endTime := calculateTimeRange(timeRange)
	fmt.Println(startTime, endTime)

	uidParsed, err := uuid.Parse(uid)
	if err != nil {
		return 0, fmt.Errorf("failed during parsing uid")
	}

	var deckIdParsed uuid.UUID
	if deckId != "" {
		deckIdParsed, err = uuid.Parse(deckId)
		if err != nil {
			return 0, fmt.Errorf("failed during parsing uid")
		}
	} else {
		deckIdParsed = uuid.UUID{}
	}

	fmt.Println("POINT4")
	avg, err := s.repo.AverageGrade(uidParsed, deckIdParsed, startTime, endTime)
	if err != nil {
		return 0, err
	}

	return avg, nil
}

func (s *Service) GetCardsReviewedCount(uid, deckId string, timeRange statsv1.TimeRange) (int32, error) {
	if timeRange == statsv1.TimeRange_TIME_RANGE_UNSPECIFIED {
		return 0, errors.New("time range must be specified")
	}

	startTime, endTime := calculateTimeRange(timeRange)

	uidParsed, err := uuid.Parse(uid)
	if err != nil {
		return 0, fmt.Errorf("failed during parsing uid")
	}

	var deckIdParsed uuid.UUID
	if deckId != "" {
		deckIdParsed, err = uuid.Parse(deckId)
		if err != nil {
			return 0, fmt.Errorf("failed during parsing uid")
		}
	} else {
		deckIdParsed = uuid.UUID{}
	}

	count, err := s.repo.CountReviewedCards(uidParsed, deckIdParsed, startTime, endTime)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Service) AddRecord(uid uuid.UUID, deckId, cardId string, createdAt time.Time, grade int) (string, error) {
	var err error
	var deckIdParsed uuid.UUID

	if deckId != "" {
		deckIdParsed, err = uuid.Parse(deckId)
		if err != nil {
			return "", fmt.Errorf("failed during parsing uid")
		}
	} else {
		deckIdParsed = uuid.UUID{}
	}

	var cardIdParsed uuid.UUID
	if cardId != "" {
		cardIdParsed, err = uuid.Parse(deckId)
		if err != nil {
			return "", fmt.Errorf("failed during parsing uid")
		}
	} else {
		cardIdParsed = uuid.UUID{}
	}

	fmt.Println("QWE", deckIdParsed, cardIdParsed)

	reviewId, err := s.repo.AddRecord(uid, deckIdParsed, cardIdParsed, createdAt, grade)
	if err != nil {
		return "", err
	}

	return reviewId, nil
}
