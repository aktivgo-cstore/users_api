package repository

import (
	"github.com/jmoiron/sqlx"
	"users_api/internal/models"
)

type UserRepository struct {
	MySqlConn *sqlx.DB
}

func NewUserRepository(mySqlConn *sqlx.DB) *UserRepository {
	return &UserRepository{
		MySqlConn: mySqlConn,
	}
}

func (ur *UserRepository) GetUsers() ([]models.User, error) {
	sql := `
		SELECT * FROM users
	`

	var users []models.User
	if err := ur.MySqlConn.Select(&users, sql); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) AddUser(user models.User) error {
	sql := `
		INSERT INTO users
		(fullName, email, hashPassword, isActivated, role) VALUE 
		(?, ?, ?, ?, ?)
	`

	args := make([]interface{}, 5)
	args = append(args, user.FullName, user.Email, user.HashPassword, user.IsActivated, user.Role)

	if _, err := ur.MySqlConn.Exec(sql, args); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) DeleteUser(id int) error {
	sql := `
		DELETE FROM users
		WHERE users.id = ?
	`

	_, err := ur.MySqlConn.Exec(sql, id)
	if err != nil {
		return err
	}

	return nil
}
