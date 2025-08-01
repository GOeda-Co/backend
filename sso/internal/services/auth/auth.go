package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	// "sso/internal/models"
	models "github.com/GOeda-Co/proto-contract/model/user"
	modelsApp "github.com/GOeda-Co/proto-contract/model/app"
	"sso/internal/storage"
	"time"

	"sso/internal/lib/logger/sl"

	// "github.com/golang-jwt/jwt/v5"
	"sso/internal/lib/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage interface {
	SaveUser(ctx context.Context, email string, hashPass []byte, name string) (uid uuid.UUID, err error)
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error)
	RegisterApp(ctx context.Context, name string, secret string) (appID int, err error)
}

// interface to get app from the storage
type AppProvider interface {
	App(ctx context.Context, appID int) (modelsApp.App, error)
}

type Auth struct {
	log         *slog.Logger
	usrStorage  UserStorage
	appProvider AppProvider
	tokenTTL    time.Duration
}

func New(
	log *slog.Logger,
	usrStorage UserStorage,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log,
		usrStorage,
		appProvider,
		tokenTTL,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, email, pass, name string) (uuid.UUID, error) {
	const op = "Auth.RegisterNewUser" //operation name for convenient logging

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
		slog.String("name", name),
	)

	log.Info("registering user")

	// bcrypt.GenerateFromPassword generates both hash with salt
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrStorage.SaveUser(ctx, email, passHash, name)
	if err != nil {
		log.Error("failed to write user to storage", sl.Err(err))
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Login checks if user with given credentials exists in the system and returns access token.
//
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error.
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string, // пароль в чистом виде, аккуратней с логами!
	appID int, // ID приложения, в котором логинится пользователь
) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
		// password либо не логируем, либо логируем в замаскированном виде
	)

	log.Info("attempting to login user")

	// Достаём пользователя из БД
	user, err := a.usrStorage.User(ctx, email)
	if err != nil {
		fmt.Printf("%v\n", err)
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}


		// a.log.Error("failed to get user", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	fmt.Println("GOT HERE", user)

	// Проверяем корректность полученного пароля
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	// Получаем информацию о приложении
	a.log.Debug("getting app information", slog.Int("app_id", appID))
	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	// Создаём токен авторизации
	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	const op = "Auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.String("user_id", userID.String()),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.usrStorage.IsAdmin(ctx, userID)
	if err != nil {
		log.Error(err.Error())
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}

func (a *Auth) RegisterApp(ctx context.Context, name string, secret string) (int, error) {
	const op = "Auth.RegisterApp"

	log := a.log.With(
		slog.String("op", op),
		slog.String("name", name),
	)

	log.Info("registering new app")

	appID, err := a.usrStorage.RegisterApp(ctx, name, secret)
	if err != nil {
		log.Error("failed to register app", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("app registered successfully", slog.Int("app_id", appID))

	return appID, nil
}
