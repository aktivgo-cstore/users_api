package models

type User struct {
	ID int `db:"id" json:"id"`
	FullName string `db:"full_name" json:"full_name"`
	Email string `db:"email" json:"email"`
	HashPassword string `db:"hash_password" json:"hash_password"`
	Role string `db:"role" json:"role"`
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
}
