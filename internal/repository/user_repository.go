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

func (ur *UserRepository) GetUsers() ([]*models.User, error) {
	sql := `
		SELECT * FROM users
	`

	var users []*models.User
	if err := ur.MySqlConn.Select(&users, sql); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) GetUser(email string) (*models.User, error) {
	var user []*models.User

	sql := `
		SELECT * FROM users
		WHERE email = ?
	`

	if err := ur.MySqlConn.Select(&user, sql, email); err != nil {
		return nil, err
	}

	if len(user) < 1 {
		return nil, nil
	}

	return user[0], nil
}

func (ur *UserRepository) SaveUser(user *models.User) error {
	sql := `
		INSERT INTO users
		(fullName, email, hashPassword, isActivated, activationLink, role, refreshToken) VALUE 
		(?, ?, ?, ?, ?, ?, ?)
	`
	var args []interface{}
	args = append(args,
		user.FullName, user.Email,
		user.HashPassword, user.IsActivated,
		user.ActivationLink, user.Role, user.RefreshToken,
	)

	if _, err := ur.MySqlConn.Exec(sql, args...); err != nil {
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
