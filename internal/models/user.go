package models

type User struct {
	ID             int         `db:"id" json:"id"`
	FullName       string      `db:"fullName" json:"fullName"`
	Email          string      `db:"email" json:"email"`
	HashPassword   string      `db:"hashPassword" json:"hashPassword"`
	IsActivated    string      `db:"isActivated" json:"isActivated"`
	ActivationLink interface{} `db:"activationLink" json:"activationLink"`
	Role           string      `db:"role" json:"role"`
	RefreshToken   string      `db:"refreshToken" json:"refreshToken"`
}
