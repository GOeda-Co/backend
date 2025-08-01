package service_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	// "sso/internal/models"
	models "github.com/GOeda-Co/proto-contract/model/user"
	modelsApp "github.com/GOeda-Co/proto-contract/model/app"
	"sso/internal/services/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserStorage struct {
	mock.Mock
}

func (m *MockUserStorage) SaveUser(ctx context.Context, email string, hashPass []byte, name string) (uuid.UUID, error) {
	args := m.Called(ctx, email, hashPass, name)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockUserStorage) User(ctx context.Context, email string) (models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserStorage) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserStorage) RegisterApp(ctx context.Context, name string, secret string) (int, error) {
	args := m.Called(ctx, name, secret)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Int(0), args.Error(1)
}

type MockAppProvider struct {
	mock.Mock
}

func (m *MockAppProvider) App(ctx context.Context, appID int) (modelsApp.App, error) {
	args := m.Called(ctx, appID)
	return args.Get(0).(modelsApp.App), args.Error(1)
}

func TestRegisterNewUser_Success(t *testing.T) {
	mockStorage := new(MockUserStorage)
	mockApps := new(MockAppProvider)
	log := slog.Default()
	service := auth.New(log, mockStorage, mockApps, time.Minute)

	ctx := context.Background()
	email := "user@example.com"
	pass := "securepass123"
	name := "Test User"
	id := uuid.New()

	mockStorage.On("SaveUser", ctx, email, mock.Anything, name).Return(id, nil)

	gotID, err := service.RegisterNewUser(ctx, email, pass, name)

	assert.NoError(t, err)
	assert.Equal(t, id, gotID)
	mockStorage.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockStorage := new(MockUserStorage)
	mockApps := new(MockAppProvider)
	log := slog.Default()
	service := auth.New(log, mockStorage, mockApps, time.Minute)

	ctx := context.Background()
	email := "user@example.com"
	password := "securepass"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{
		ID:       uuid.New(),
		Email:    email,
		PassHash: hashed,
	}

	app := modelsApp.App{ID: 1, Name: "TestApp"}

	mockStorage.On("User", ctx, email).Return(user, nil)
	mockApps.On("App", ctx, 1).Return(app, nil)

	token, err := service.Login(ctx, email, password, 1)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockStorage := new(MockUserStorage)
	mockApps := new(MockAppProvider)
	log := slog.Default()
	service := auth.New(log, mockStorage, mockApps, time.Minute)

	ctx := context.Background()
	email := "user@example.com"
	wrongPass := "wrong"
	hashed, _ := bcrypt.GenerateFromPassword([]byte("realpass"), bcrypt.DefaultCost)
	user := models.User{
		Email:    email,
		PassHash: hashed,
	}

	mockStorage.On("User", ctx, email).Return(user, nil)

	_, err := service.Login(ctx, email, wrongPass, 1)

	assert.ErrorIs(t, err, auth.ErrInvalidCredentials)
}
