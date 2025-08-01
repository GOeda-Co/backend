package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	// "sso/internal/models"
	"sso/internal/storage"
	"sso/migrations"

	modelsApp "github.com/GOeda-Co/proto-contract/model/app"
	models "github.com/GOeda-Co/proto-contract/model/user"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	DB  *gorm.DB
	log *slog.Logger
}

func New(connString string, log *slog.Logger) (*Storage, error) {
	const (
		op         = "storage.postgresql.New"
		maxRetries = 10
		retryDelay = 4 * time.Second
	)

	var DB *gorm.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(connString), &gorm.Config{})
		if err == nil {
			// Try pinging to ensure DB is really ready
			sqlDB, errPing := DB.DB()
			if errPing == nil && sqlDB.Ping() == nil {
				break // success
			}
		}
		log.Info("waiting for database...",
			"operation", op,
			"attempt", fmt.Sprintf("%d/%d", i+1, maxRetries))
		time.Sleep(retryDelay)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect after %d attempts: %w", op, maxRetries, err)
	}

	migrations.MigrateToLatest(DB, log)

	// AutoMigrate your models
	if err := DB.AutoMigrate(&models.User{}, &modelsApp.App{}); err != nil {
		return nil, fmt.Errorf("%s: migration error: %w", op, err)
	}

	return &Storage{DB: DB, log: log}, nil
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
	user := models.User{
		Email:    email,
		PassHash: hashPass,
		Name:     name,
	}
	err = s.DB.Create(&user).Error
	return user.ID, fmt.Errorf("%s: %w", op, err)
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "Storage.postgresql.User"
	var user models.User
	if err := s.DB.Where("email = ?", email).Find(&user).Error; err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) App(ctx context.Context, appID int) (modelsApp.App, error) {
	const op = "Storage.postgresql.App"
	var app modelsApp.App
	if err := s.DB.Where("ID = ?", appID).Find(&app).Error; err != nil {
		return modelsApp.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	const op = "Storage.postgresql.IsAdmin"
	var user models.User
	err := s.DB.WithContext(ctx).First(&user, "ID = ?", userID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("%s: %v", op, storage.ErrUserNotFound)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return false, fmt.Errorf("%s: %v", op, err)
	}
	return user.IsAdmin, nil
}

func (s *Storage) RegisterApp(ctx context.Context, name string, secret string) (appID int, err error) {
	const op = "Storage.postgresql.RegisterApp"
	app := modelsApp.App{
		Name:   name,
		Secret: secret,
	}
	err = s.DB.Create(&app).Error
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return app.ID, nil
}