package service

import (
	"github.com/golang-jwt/jwt"
	"os"
	"time"
	"users_api/internal/dto"
	"users_api/internal/types"
)

var (
	accessKey  = os.Getenv("JWT_ACCESS_SECRET")
	refreshKey = os.Getenv("JWT_REFRESH_SECRET")
)

func GenerateTokens(tokenData *dto.TokenData) (*types.Tokens, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    tokenData.ID,
		"email": tokenData.Email,
		"role":  tokenData.Role,
		"exp":   time.Now().Add(time.Minute * 30).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(accessKey))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    tokenData.ID,
		"email": tokenData.Email,
		"role":  tokenData.Role,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(refreshKey))
	if err != nil {
		return nil, err
	}

	tokens := &types.Tokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}

	return tokens, nil
}
