package postgres_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tomatoCoderq/card/internal/config"

	// "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/tomatoCoderq/card/internal/repository/postgresql"
	"github.com/tomatoCoderq/card/pkg/model"
)

var repo *postgresql.Repository
var db *gorm.DB

func TestMain(m *testing.M) {
	cfg := config.MustLoad()
	log := slog.Default()

	dsn := cfg.ConnectionString

	repo = postgresql.New(dsn, log)

	code := m.Run()

	// Clean up after tests
	os.Exit(code)
}

func TestCardRepo_CRUD(t *testing.T) {
	userId, _ := uuid.Parse("e1d3ce6e-1e99-4724-b57a-45e0c5d9dd08")
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

	cards, err := repo.ReadAllCards(userId)
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
