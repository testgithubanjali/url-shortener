package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/testgithubanjali/url-shortener/internal/config"
)

func GenerateJWT(userID string) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}
