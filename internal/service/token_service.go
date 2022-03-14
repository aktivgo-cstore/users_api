package service

import (
	"github.com/golang-jwt/jwt"
	"os"
	"time"
	"users_api/internal/dto"
)

var (
	accessKey = os.Getenv("JWT_ACCESS_SECRET")
)

func GenerateToken(tokenData *dto.TokenData) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    tokenData.ID,
		"email": tokenData.Email,
		"role":  tokenData.Role,
		"exp":   time.Now().Add(time.Minute * 30).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(accessKey))
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}
