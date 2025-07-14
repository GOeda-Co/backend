package postgres_test

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tomatoCoderq/deck/internal/config"
	"github.com/tomatoCoderq/deck/internal/repository/postgresql"
	"github.com/tomatoCoderq/deck/pkg/model"
)

var testRepo *postgresql.Repository

func TestMain(m *testing.M) {

	cfg := config.MustLoad()
	log := slog.Default()

	testRepo = postgresql.New(cfg.ConnectionString, log)

	code := m.Run()
	os.Exit(code)
}

func TestAddAndReadDeck(t *testing.T) {
	t.Helper()
	deckId := uuid.New()
	userId := uuid.New()

	deck := &model.Deck{
		DeckId:    deckId,
		Name:      "Test Deck",
		CreatedBy: userId,
		CreatedAt: time.Now(),
	}

	err := testRepo.AddDeck(deck)
	assert.NoError(t, err)

	
	readDeck, err := testRepo.ReadDeck(deckId)
	assert.NoError(t, err)
	assert.Equal(t, "Test Deck", readDeck.Name)
	assert.Equal(t, userId, readDeck.CreatedBy)
}

func TestReadAllDecksOfUser(t *testing.T) {
	userId := uuid.New()

	deck := &model.Deck{
		DeckId:    uuid.New(),
		Name:      "User Deck",
		CreatedBy: userId,
		CreatedAt: time.Now(),
	}

	_ = testRepo.AddDeck(deck)

	decks, err := testRepo.ReadAllDecksOfUser(userId)
	assert.NoError(t, err)
	assert.True(t, len(decks) > 0)
	assert.Equal(t, userId, decks[0].CreatedBy)
}

func TestDeleteDeck(t *testing.T) {
	deck := &model.Deck{
		DeckId:    uuid.New(),
		Name:      "Delete Me",
		CreatedBy: uuid.New(),
		CreatedAt: time.Now(),
	}

	_ = testRepo.AddDeck(deck)

	err := testRepo.DeleteDeck(deck.DeckId)
	assert.NoError(t, err)

	_, err = testRepo.ReadDeck(deck.DeckId)
	assert.Error(t, err)
}