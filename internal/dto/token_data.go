package dto

type TokenData struct {
	ID    int64
	Email string
	Role  string
}

func NewTokenData(id int64, email string, role string) *TokenData {
	return &TokenData{
		ID:    id,
		Email: email,
		Role:  role,
	}
}
