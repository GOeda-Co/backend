package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"sso/internal/models"
	"sso/internal/storage"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func New(connString string, log *slog.Logger) (*Storage, error) {
	const (
		op        = "storage.postgresql.New"
		maxRetries = 10
		retryDelay = 4 * time.Second
	)

	var db *gorm.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(connString), &gorm.Config{})
		if err == nil {
			// Try pinging to ensure DB is really ready
			sqlDB, errPing := db.DB()
			if errPing == nil && sqlDB.Ping() == nil {
				break // success
			}
		}

		log.Info("[%s] waiting for database... (%d/%d)", op, i+1, maxRetries)
		time.Sleep(retryDelay)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect after %d attempts: %w", op, maxRetries, err)
	}

	// AutoMigrate your models
	if err := db.AutoMigrate(&models.User{}, &models.App{}); err != nil {
		return nil, fmt.Errorf("%s: migration error: %w", op, err)
	}

	return &Storage{db: db}, nil
}

/*
FOR TESTS
type UserStorage interface {
	SaveUser(ctx context.Context, email string, hashPass []byte) (uid int64, err error)
	User(ctx context.Context, email string) (models.User, error)
}

// interface to get app from the storage
type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}
*/

func (s *Storage) SaveUser(ctx context.Context, email string, hashPass []byte, name string) (uid uuid.UUID, err error) {
	const op = "Storage.postgresql.SaveUser"
	//TODO: Cancel out this to service layer
	user := models.User{
		Email:    email,
		PassHash: hashPass,
		Name: name,
	}
	err = s.db.Create(&user).Error
	return user.ID, err
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "Storage.postgresql.User"
	var user models.User
	err := s.db.Where("email = ?", email).Find(&user).Error
	return user, err
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "Storage.postgresql.App"
	var app models.App
	err := s.db.Where("ID = ?", appID).Find(&app).Error
	return app, err
}

func (s *Storage) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	const op = "Storage.postgresql.IsAdmin"
	var user models.User
	err := s.db.WithContext(ctx).First(&user, "ID = ?", userID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("%s: %v", op, storage.ErrUserNotFound)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return false, fmt.Errorf("%s: %v", op, err)
	}
	return user.IsAdmin, nil
}

/*

func (cr Repository) AddCard(card *model.Card) error {
	return cr.db.Create(card).Error
}

func (cr Repository) ReadAllCards(userId uuid.UUID) ([]model.Card, error) {
	var cards []model.Card
	err := cr.db.Where("expires_at < ?", time.Now()).Where("created_by = ?", userId).Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, err
}

func (cr Repository) ReadCard(cardId uuid.UUID) (*model.Card, error) {
	var card model.Card
	err := cr.db.Where("card_id = ?", cardId).Find(&card).Error
	return &card, err
}

func (cr Repository) UpdateCard(card *model.Card, cardUpdate *schemes.UpdateCardScheme) (*model.Card, error) {
	updateCardFields(card, cardUpdate)
	return card, cr.db.Updates(card).Error
}

func (cr Repository) PureUpdate(card *model.Card) error {
	return cr.db.Updates(card).Error
}

func (cr Repository) DeleteCard(cardId uuid.UUID) error {
	err := cr.db.Delete(&model.Card{}, "card_id = ?", cardId).Error
	return err
}

*/
