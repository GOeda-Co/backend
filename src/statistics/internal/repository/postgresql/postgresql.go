package postgresql

import (
	"log"
	"repeatro/src/statistics/pkg/model"

	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Need to cancel out to service layer
type Repository struct {
	db *gorm.DB
}

func NewPostgresRepo(config *viper.Viper, newLogger logger.Interface) *Repository {
	db, err := gorm.Open(postgres.Open(config.GetString("database.connection_string")), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatalf("Error during opening database")
	}

	db.AutoMigrate(&model.Stat{})

	return &Repository{db: db}
}

func (r *Repository) Add(result *model.Stat) error {
	return r.db.Create(result).Error
}

func (r *Repository) Delete(resultId uuid.UUID) error {
	return r.db.Delete(&model.Stat{}, "id = ?", resultId).Error
}

func (r *Repository) GetAllGradesForPeriod(dtStart, dtEnd time.Time, userId uuid.UUID) ([]int, error) {
	var grades []int
	err := r.db.
		Model(&model.Stat{}).
		Where("created_at BETWEEN ? AND ?", dtStart, dtEnd).
		Pluck("grade", &grades).Error
	return grades, err
}

func (r *Repository) GetLearnedCardsForPeriod(dtStart, dtEnd time.Time, userId uuid.UUID) ([]uuid.UUID, error) {
	var cardIDs []uuid.UUID
	err := r.db.
		Model(&model.Stat{}).
		Select("DISTINCT card_id").
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userId, dtStart, dtEnd).
		Pluck("card_id", &cardIDs).Error
	return cardIDs, err
}



