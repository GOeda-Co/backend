package postgresql

import (
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/tomatoCoderq/stats/pkg/model/review"

	// "github.com/tomatoCoderq/card/pkg/scheme"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// type Repository interface {
// 	AverageGrade(uid, cardId string, startTime, endTime time.Time) (float64, error)
// 	CountReviewedCards(uid, cardId string, startTime, endTime time.Time) (int32, error)
// GetCardsLearnedCount(uid, cardId string, startTime, endTime time.Time) (int32, error)
// }

func FormExec(uid, deckId uuid.UUID, startTime, endTime time.Time, cr Repository) *gorm.DB {
	var exec *gorm.DB

	exec = cr.db.Table("results")

	if reflect.DeepEqual(uid, uuid.UUID{}) && reflect.DeepEqual(deckId, uuid.UUID{}) {
		exec = exec.Where("created_at BETWEEN ? AND ?", startTime, endTime)
	} else if reflect.DeepEqual(uid, uuid.UUID{}) {
		exec = exec.Where("deck_id = ?", deckId).Where("created_at BETWEEN ? AND ?", startTime, endTime)
	} else if reflect.DeepEqual(deckId, uuid.UUID{}) {
		fmt.Println(">???", startTime, endTime, uid)
		exec = exec.Where("user_id = ?", uid).Where("created_at BETWEEN ? AND ?", startTime, endTime)
	} else {
		exec = exec.Where("user_id = ?", uid).Where("deck_id = ?", deckId).Where("created_at BETWEEN ? AND ?", startTime, endTime)
	}
	return exec
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

	db.AutoMigrate(&model.Review{})

	return &Repository{db: db}
}

// TODO: Divide into separated functions
func (cr Repository) AverageGrade(uid, deckId uuid.UUID, startTime, endTime time.Time) (float64, error) {
	exec := FormExec(uid, deckId, startTime, endTime, cr)

	var reviews []model.Review
	if err := exec.Find(&reviews).Error; err != nil {
		return 0, err
	}

	if len(reviews) == 0 {
		return 0, nil
	}

	average := func(reviews []model.Review) float64 {
		var sum int
		for _, review := range reviews {
			sum += int(review.Grade)
		}
		return float64(sum) / float64(len(reviews))
	}

	return average(reviews), nil
}

func (cr Repository) CountReviewedCards(uid, deckId uuid.UUID, startTime, endTime time.Time) (int32, error) {
	exec := FormExec(uid, deckId, startTime, endTime, cr)

	var reviews []model.Review
	if err := exec.Find(&reviews).Error; err != nil {
		return 0, err
	}

	return int32(len(reviews)), nil
}

func (cr Repository) AddRecord(uid, deckId, cardId uuid.UUID, createdAt time.Time, grade int) (string, error) {
	// var review model.Review
	review := model.Review{
		UserID:    uid,
		DeckId:    deckId,
		CardID:    cardId,
		CreatedAt: createdAt,
		Grade:     int32(grade),
	}
	if err := cr.db.Create(&review).Error; err != nil {
		return "", err
	}

	return review.ResultId.String(), nil
}
