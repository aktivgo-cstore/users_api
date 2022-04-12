package service

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"os"
	"time"
	"users_api/internal/constants"
	"users_api/internal/dto"
	"users_api/internal/errors"
)

var (
	key = os.Getenv("JWT_ACCESS_SECRET")
)

type Claims struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	IsActivated int    `json:"is_activated"`
	Role        string `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(tokenData *dto.TokenData) (string, error) {
	expirationTime := time.Now().Add(500 * time.Hour)
	claims := &Claims{
		ID:          tokenData.ID,
		Email:       tokenData.Email,
		IsActivated: tokenData.IsActivated,
		Role:        tokenData.Role,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetToken(header http.Header) (string, *errors.ApiError) {
	token := header.Get("token")
	if token == "" {
		return "", errors.UnauthorizedError()
	}

	return token, nil
}

func CheckAccess(tokenStr string) *errors.ApiError {
	tokenData, err := decodeToken(tokenStr)
	if err != nil {
		log.Println("unable to decode token:", err.Error())
		return errors.UnauthorizedError()
	}

	if tokenData.Role != constants.RoleAdmin {
		return errors.ForbiddenError()
	}

	return nil
}

func decodeToken(tokenStr string) (*dto.TokenData, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}

	id, ok := claims["id"]
	if !ok {
		return nil, fmt.Errorf("id not found")
	}

	email, ok := claims["email"]
	if !ok {
		return nil, fmt.Errorf("email not found")
	}

	isActivated, ok := claims["is_activated"]
	if !ok {
		return nil, fmt.Errorf("isActivated not found")
	}

	role, ok := claims["role"]
	if !ok {
		return nil, fmt.Errorf("role not found")
	}

	return &dto.TokenData{
		ID:          int(id.(float64)),
		Email:       email.(string),
		IsActivated: int(isActivated.(float64)),
		Role:        role.(string),
	}, nil
}
