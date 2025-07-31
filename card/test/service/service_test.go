package services_test

import (
	"context"
	"testing"
	"time"

	// "time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"

	"log/slog"

	services "github.com/tomatoCoderq/card/internal/services/card"
	"github.com/tomatoCoderq/card/pkg/model"
	schemes "github.com/tomatoCoderq/card/pkg/scheme"
)

type MockCardRepo struct {
	mock.Mock
}

// Mock Stats Client
type MockStatsClient struct {
	mock.Mock
}

func (m *MockStatsClient) AddRecord(ctx context.Context, deckId, cardId string, grade int) (string, error) {
	args := m.Called(ctx, deckId, cardId, grade)
	return args.String(0), args.Error(1)
}

// ReadAllCardsByUser implements services.CardRepository.
func (m *MockCardRepo) ReadAllCardsByUser(userId uuid.UUID) ([]model.Card, error) {
	panic("unimplemented")
}

func (m *MockCardRepo) AddCard(card *model.Card) error {
	args := m.Called(card)
	return args.Error(0)
}

func (m *MockCardRepo) ReadAllCards(userId uuid.UUID) ([]model.Card, error) {
	args := m.Called(userId)
	return args.Get(0).([]model.Card), args.Error(1)
}

func (m *MockCardRepo) ReadCard(cardId uuid.UUID) (*model.Card, error) {
	args := m.Called(cardId)
	return args.Get(0).(*model.Card), args.Error(1)
}

func (m *MockCardRepo) PureUpdate(card *model.Card) error {
	args := m.Called(card)
	return args.Error(0)
}

func (m *MockCardRepo) UpdateCard(card *model.Card, cardUpdate *schemes.UpdateCardScheme) (*model.Card, error) {
	args := m.Called(card, cardUpdate)
	return args.Get(0).(*model.Card), args.Error(1)
}

func (m *MockCardRepo) DeleteCard(cardId uuid.UUID) error {
	args := m.Called(cardId)
	return args.Error(0)
}

func (m *MockCardRepo) SearchAllPublicCards() ([]model.Card, error) {
	args := m.Called()
	return args.Get(0).([]model.Card), args.Error(1)
}

func (m *MockCardRepo) SearchUserPublicCards(userId uuid.UUID) ([]model.Card, error) {
	args := m.Called(userId)
	return args.Get(0).([]model.Card), args.Error(1)
}

func TestAddCard(t *testing.T) {
	mockRepo := new(MockCardRepo)
	logger := slog.Default()
	service := services.New(logger, mockRepo, nil)

	card := &model.Card{
		CardId:      uuid.New(),
		CreatedBy:   uuid.New(),
		Word:        "front",
		Translation: "back",
	}

	mockRepo.On("AddCard", card).Return(nil)

	result, err := service.AddCard(card)

	assert.NoError(t, err)
	assert.Equal(t, card, result)
	mockRepo.AssertExpectations(t)
}

func TestReadAllCards(t *testing.T) {
	mockRepo := new(MockCardRepo)
	logger := slog.Default()
	service := services.New(logger, mockRepo, nil)

	userId := uuid.New()
	expectedCards := []model.Card{{Word: "A"}, {Translation: "B"}}
	mockRepo.On("ReadAllCards", userId).Return(expectedCards, nil)

	cards, err := service.ReadAllCards(userId)

	assert.NoError(t, err)
	assert.Len(t, cards, 2)
	mockRepo.AssertExpectations(t)
}

func TestUpdateCard(t *testing.T) {
	mockRepo := new(MockCardRepo)
	logger := slog.Default()
	service := services.New(logger, mockRepo, nil)

	cardId := uuid.New()
	userId := uuid.New()
	card := &model.Card{CardId: cardId, CreatedBy: userId, Word: "old"}
	updatedCard := &model.Card{CardId: cardId, CreatedBy: userId, Word: "new"}
	update := &schemes.UpdateCardScheme{Word: "new"}

	mockRepo.On("ReadCard", cardId).Return(card, nil)
	mockRepo.On("UpdateCard", card, update).Return(updatedCard, nil)

	result, err := service.UpdateCard(cardId, update, userId)

	assert.NoError(t, err)
	assert.Equal(t, updatedCard, result)
	mockRepo.AssertExpectations(t)
}

func TestDeleteCard(t *testing.T) {
	mockRepo := new(MockCardRepo)
	logger := slog.Default()
	service := services.New(logger, mockRepo, nil)

	cardId := uuid.New()
	userId := uuid.New()
	card := &model.Card{CardId: cardId, CreatedBy: userId}

	mockRepo.On("ReadCard", cardId).Return(card, nil)
	mockRepo.On("DeleteCard", cardId).Return(nil)

	err := service.DeleteCard(cardId, userId)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAddAnswers_ValidGradeAndOwner(t *testing.T) {
	mockRepo := new(MockCardRepo)
	mockStatsClient := new(MockStatsClient)
	logger := slog.Default()
	service := services.New(logger, mockRepo, mockStatsClient)

	userId := uuid.New()
	cardId := uuid.New()
	deckId := uuid.New()
	card := &model.Card{
		CardId:           cardId,
		CreatedBy:        userId,
		DeckID:           deckId,
		ExpiresAt:        time.Now().Add(-time.Hour),
		Interval:         1,
		Easiness:         2.5,
		RepetitionNumber: 1,
	}

	answers := []schemes.AnswerScheme{
		{CardId: cardId, Grade: 4},
	}

	// Create a JWT token-like string with proper format
	mockJWT := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiIxMjM0NTY3OC05YWJjLWRlZjAtMTIzNC01Njc4OWFiY2RlZjAiLCJlbWFpbCI6InRlc3RAdGVzdC5jb20iLCJleHAiOjk5OTk5OTk5OTksIm5hbWUiOiJUZXN0IFVzZXIifQ.test_signature"

	mockRepo.On("ReadCard", cardId).Return(card, nil)
	mockRepo.On("PureUpdate", mock.AnythingOfType("*model.Card")).Return(nil)
	mockStatsClient.On("AddRecord", mock.Anything, deckId.String(), cardId.String(), 4).Return("review-id-123", nil)

	// Create context with proper JWT authorization metadata
	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		"authorization": mockJWT,
		"user_id":       userId.String(),
	}))
	
	err := service.AddAnswers(ctx, userId, answers)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockStatsClient.AssertExpectations(t)
}

func TestSearchAllPublicCards(t *testing.T) {
	mockRepo := new(MockCardRepo)
	logger := slog.Default()
	service := services.New(logger, mockRepo, nil)

	expectedCards := []model.Card{{Word: "A", Translation: "B"}, {Word: "C", Translation: "D"}}
	mockRepo.On("SearchAllPublicCards").Return(expectedCards, nil)

	cards, err := service.SearchAllPublicCards()

	assert.NoError(t, err)
	assert.Len(t, cards, 2)
	mockRepo.AssertExpectations(t)
}

func TestSearchUserPublicCards(t *testing.T) {
	mockRepo := new(MockCardRepo)
	logger := slog.Default()
	service := services.New(logger, mockRepo, nil)

	userId := uuid.New()
	expectedCards := []model.Card{{Word: "A", Translation: "B"}, {Word: "C", Translation: "D"}}
	mockRepo.On("SearchUserPublicCards", userId).Return(expectedCards, nil)

	cards, err := service.SearchUserPublicCards(userId.String())

	assert.NoError(t, err)
	assert.Len(t, cards, 2)
	mockRepo.AssertExpectations(t)
}