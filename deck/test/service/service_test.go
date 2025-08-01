package services_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	services "github.com/tomatoCoderq/deck/internal/services/deck"
	// schemes "github.com/GOeda-Co/proto-contract/scheme/deck"
	"github.com/GOeda-Co/proto-contract/model/deck"
	modelCard "github.com/GOeda-Co/proto-contract/model/card"
)

type MockDeckRepository struct {
	mock.Mock
}

func (m *MockDeckRepository) AddDeck(deck *model.Deck) error {
	args := m.Called(deck)
	return args.Error(0)
}

func (m *MockDeckRepository) ReadAllDecksOfUser(userId uuid.UUID) ([]model.Deck, error) {
	args := m.Called(userId)
	return args.Get(0).([]model.Deck), args.Error(1)
}

func (m *MockDeckRepository) ReadAllDecks() ([]model.Deck, error) {
	args := m.Called()
	return args.Get(0).([]model.Deck), args.Error(1)
}

func (m *MockDeckRepository) ReadDeck(deckId uuid.UUID) (*model.Deck, error) {
	args := m.Called(deckId)
	return args.Get(0).(*model.Deck), args.Error(1)
}

func (m *MockDeckRepository) DeleteDeck(deckId uuid.UUID) error {
	args := m.Called(deckId)
	return args.Error(0)
}

func (m *MockDeckRepository) SearchAllPublicDecks() ([]model.Deck, error) {
	args := m.Called()
	return args.Get(0).([]model.Deck), args.Error(1)
}
func (m *MockDeckRepository) SearchUserPublicDecks(userId uuid.UUID) ([]model.Deck, error) {
	args := m.Called(userId)
	return args.Get(0).([]model.Deck), args.Error(1)
}

func (m *MockDeckRepository) AddCardToDeck(cardId uuid.UUID, deckId uuid.UUID) error {
	args := m.Called(cardId, deckId)
	return args.Error(0)
}

func (m *MockDeckRepository) FindAllCardsInDeck(deckId uuid.UUID) ([]modelCard.Card, error) {
	args := m.Called(deckId)
	return args.Get(0).([]modelCard.Card), args.Error(1)
}

func TestAddDeck(t *testing.T) {
	mockRepo := new(MockDeckRepository)
	service := services.New(nil, mockRepo)

	deck := &model.Deck{
		DeckId:    uuid.New(),
		Name:      "Test Deck",
		CreatedBy: uuid.New(),
		CreatedAt: time.Now(),
	}

	mockRepo.On("AddDeck", deck).Return(nil)

	result, err := service.AddDeck(deck)
	assert.NoError(t, err)
	assert.Equal(t, deck, result)
	mockRepo.AssertExpectations(t)
}

func TestReadDeck_Authorized(t *testing.T) {
	mockRepo := new(MockDeckRepository)
	service := services.New(nil, mockRepo)

	userId := uuid.New()
	deck := &model.Deck{
		DeckId:    uuid.New(),
		Name:      "Private Deck",
		CreatedBy: userId,
	}

	mockRepo.On("ReadDeck", deck.DeckId).Return(deck, nil)

	result, err := service.ReadDeck(deck.DeckId, userId)
	assert.NoError(t, err)
	assert.Equal(t, deck, result)
}

func TestReadDeck_Unauthorized(t *testing.T) {
	mockRepo := new(MockDeckRepository)
	service := services.New(nil, mockRepo)

	ownerId := uuid.New()
	requesterId := uuid.New()
	deck := &model.Deck{
		DeckId:    uuid.New(),
		Name:      "Private Deck",
		CreatedBy: ownerId,
	}

	mockRepo.On("ReadDeck", deck.DeckId).Return(deck, nil)

	result, err := service.ReadDeck(deck.DeckId, requesterId)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, services.ErrUnauthorized, err)
}

func TestDeleteDeck_Success(t *testing.T) {
	mockRepo := new(MockDeckRepository)
	service := services.New(nil, mockRepo)

	userId := uuid.New()
	deckId := uuid.New()
	deck := &model.Deck{
		DeckId:    deckId,
		CreatedBy: userId,
	}

	mockRepo.On("ReadDeck", deckId).Return(deck, nil)
	mockRepo.On("DeleteDeck", deckId).Return(nil)

	err := service.DeleteDeck(deckId, userId)
	assert.NoError(t, err)
}

func TestAddCardToDeck_Unauthorized(t *testing.T) {
	mockRepo := new(MockDeckRepository)
	service := services.New(nil, mockRepo)

	ownerId := uuid.New()
	requesterId := uuid.New()
	deckId := uuid.New()
	cardId := uuid.New()

	deck := &model.Deck{
		DeckId:    deckId,
		CreatedBy: ownerId,
	}

	mockRepo.On("ReadDeck", deckId).Return(deck, nil)

	err := service.AddCardToDeck(cardId, deckId, requesterId)
	assert.Error(t, err)
	assert.Equal(t, services.ErrUnauthorized, err)
}
