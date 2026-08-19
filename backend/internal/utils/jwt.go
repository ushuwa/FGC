package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/jj.jobo/FGC/internal/config"
)

func GenerateJWT(id uint, username string, role string) (string, error) {

	claims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(config.App.JWTSecret))
}
