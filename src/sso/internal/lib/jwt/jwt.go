package jwt

import (
	"fmt"
	"time"

	"github.com/HAHLIK/AuthService/sso/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString(app.Sercret)
	if err != nil {
		return "", fmt.Errorf("%s : %w", "can't create user token", err)
	}

	return tokenString, nil
}
