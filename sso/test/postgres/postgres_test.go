package postgresql

import (
	"context"
	"fmt"
	"math/rand/v2"

	"log/slog"
	"os"
	"testing"

	"sso/config"

	modelsApp "github.com/GOeda-Co/proto-contract/model/app"
	"sso/internal/storage/postgresql"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var repo *postgresql.Storage

func TestMain(m *testing.M) {
	cfg := config.MustLoad()
	log := slog.Default()

	dsn := cfg.ConnectionString

	repo, _ = postgresql.New(dsn, log)

	code := m.Run()

	os.Exit(code)
}

func TestSaveAndGetUser(t *testing.T) {
	ctx := context.Background()

	email := "test@example.com"
	name := "Test User"
	pass := []byte("hashedpassword123")

	uid, err := repo.SaveUser(ctx, email, pass, name)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	// fmt.Println("hererer", err.Error())
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, uid)

	user, err := repo.User(ctx, email)
	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, pass, user.PassHash)
}

func TestApp(t *testing.T) {
	ctx := context.Background()

	id := rand.IntN(1000) + 1000 + rand.IntN(rand.IntN(1000))

	app := modelsApp.App{
		ID:     id,
		Name:   "Repeatro",
		Secret: "some secret",
	}
	err := repo.DB.Create(&app).Error
	assert.NoError(t, err)

	foundApp, err := repo.App(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, "Repeatro", foundApp.Name)
}