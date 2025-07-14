package services_test

import (
	"testing"
	"time"
	// "time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"log/slog"

	"github.com/tomatoCoderq/card/internal/services/card"
	"github.com/tomatoCoderq/card/pkg/model"
	schemes "github.com/tomatoCoderq/card/pkg/scheme"
)

type MockCardRepo struct {
	mock.Mock
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

func TestAddCard(t *testing.T) {
	mockRepo := new(MockCardRepo)
	logger := slog.Default()
	service := services.New(logger, mockRepo)

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
	service := services.New(logger, mockRepo)

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
	service := services.New(logger, mockRepo)

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
	service := services.New(logger, mockRepo)

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
	logger := slog.Default()
	service := services.New(logger, mockRepo)

	userId := uuid.New()
	cardId := uuid.New()
	card := &model.Card{
		CardId:           cardId,
		CreatedBy:        userId,
		ExpiresAt:        time.Now().Add(-time.Hour),
		Interval:         1,
		Easiness:         2.5,
		RepetitionNumber: 1,
	}

	answers := []schemes.AnswerScheme{
		{CardId: cardId, Grade: 4},
	}

	mockRepo.On("ReadCard", cardId).Return(card, nil)
	mockRepo.On("PureUpdate", mock.AnythingOfType("*model.Card")).Return(nil)

	err := service.AddAnswers(userId, answers)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}