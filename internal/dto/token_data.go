package dto

type TokenData struct {
	ID          int
	Email       string
	IsActivated int
	Role        string
}

func NewTokenData(id int, email string, isActivated int, role string) *TokenData {
	return &TokenData{
		ID:          id,
		Email:       email,
		IsActivated: isActivated,
		Role:        role,
	}
}
