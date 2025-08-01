package postgres_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/tomatoCoderq/card/internal/config"

	"github.com/GOeda-Co/proto-contract/model/card"
	"github.com/tomatoCoderq/card/internal/repository/postgresql"
)

var repo *postgresql.Repository

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load("../../../.env")
	_ = godotenv.Load(".env")

	log := slog.Default()

	// Try test-specific config paths
	testConfigPaths := []string{
		"../../config/local.yaml",
		"../../config/config.yaml",
		"../../../config/local.yaml",
	}

	for _, path := range testConfigPaths {
		if _, err := os.Stat(path); err == nil {
			if err := os.Setenv("CONFIG_PATH", path); err != nil {
				log.Error("Error setting CONFIG_PATH", "error", err)
			}
			break
		}
	}

	cfg := config.MustLoad()

	dsn := cfg.ConnectionString

	repo = postgresql.New(dsn, log)

	code := m.Run()

	// Clean up after tests
	os.Exit(code)
}

func TestCardRepo_CRUD(t *testing.T) {
	userId, _ := uuid.Parse("d9f5e420-9303-45b8-90ba-f036ef0aea44")
	cardId := uuid.New()

	card := &model.Card{
		CardId:           cardId,
		CreatedBy:        userId,
		Word:             "Front side",
		Translation:      "Back side",
		Easiness:         2.5,
		Interval:         1,
		RepetitionNumber: 0,
		ExpiresAt:        time.Now().Add(-time.Hour),
		CreatedAt:        time.Now(),
	}

	err := repo.AddCard(card)
	assert.NoError(t, err)

	readCard, err := repo.ReadCard(card.CardId)
	assert.NoError(t, err)
	assert.Equal(t, "Front side", readCard.Word)

	cards, err := repo.ReadAllOwnCardsToLearn(userId)
	assert.NoError(t, err)
	fmt.Println(len(cards) >= 1)
	assert.True(t, len(cards) >= 1)

	card.Word = "Updated front"
	err = repo.PureUpdate(card)
	assert.NoError(t, err)

	updated, _ := repo.ReadCard(cardId)
	assert.Equal(t, "Updated front", updated.Word)

	err = repo.DeleteCard(cardId)
	assert.NoError(t, err)

	deleted, _ := repo.ReadCard(cardId)
	assert.Equal(t, &model.Card{}, deleted)
}

func TestSearchAllPublicCards(t *testing.T) {
	/*
		Search all public cards available in the repository
		Here I assume that there at least one public card exists
		so the result should not be empty
	*/
	results, err := repo.SearchAllPublicCards() // Search all cards that are public
	assert.NoError(t, err)
	assert.NotEmpty(t, results)
}

func TestSearchUserPublicCards(t *testing.T) {
	userId, _ := uuid.Parse("e1d3ce6e-1e99-4724-b57a-45e0c5d9dd08")
	/*
		Search all public cards for a specific user
		Here I use a user ID that is expected to not have any public cards
		so the result should be empty
	*/
	results, err := repo.SearchUserPublicCards(userId)
	assert.NoError(t, err)
	assert.Empty(t, results)
}
