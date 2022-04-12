package models

type User struct {
	ID             int64       `db:"id" json:"id"`
	FullName       string      `db:"fullName" json:"fullName"`
	Email          string      `db:"email" json:"email"`
	HashPassword   string      `db:"hashPassword" json:"hashPassword"`
	IsActivated    int         `db:"isActivated" json:"isActivated"`
	ActivationLink interface{} `db:"activationLink" json:"activationLink"`
	Role           string      `db:"role" json:"role"`
	Token          string      `db:"token" json:"token"`
}
