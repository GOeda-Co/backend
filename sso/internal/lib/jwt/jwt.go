package jwt

import (
	// "sso/internal/models"
	models "github.com/GOeda-Co/proto-contract/model/user"
	modelsApp "github.com/GOeda-Co/proto-contract/model/app"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, app modelsApp.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Добавляем в токен всю необходимую информацию
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["id"] = user.ID
	claims["admin"] = user.IsAdmin
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	// Подписываем токен, используя секретный ключ приложения
	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
