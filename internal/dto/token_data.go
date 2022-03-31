package dto

type TokenData struct {
	ID          int64
	Email       string
	IsActivated int
	Role        string
}

func NewTokenData(id int64, email string, isActivated int, role string) *TokenData {
	return &TokenData{
		ID:          id,
		Email:       email,
		IsActivated: isActivated,
		Role:        role,
	}
}
