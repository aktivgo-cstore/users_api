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

func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
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

func (ur *UserRepository) GetUserByID(id int64) (*models.User, error) {
	var user []*models.User

	sql := `
		SELECT * FROM users
		WHERE id = ?
	`

	if err := ur.MySqlConn.Select(&user, sql, id); err != nil {
		return nil, err
	}

	if len(user) < 1 {
		return nil, nil
	}

	return user[0], nil
}

func (ur *UserRepository) GetUserByActivationLink(activationLink string) (*models.User, error) {
	var user []*models.User

	sql := `
		SELECT * FROM users
		WHERE activationLink = ?
	`

	if err := ur.MySqlConn.Select(&user, sql, activationLink); err != nil {
		return nil, err
	}

	if len(user) < 1 {
		return nil, nil
	}

	return user[0], nil
}

func (ur *UserRepository) SaveUser(user *models.User) (int64, error) {
	sql := `
		INSERT INTO users
		(fullName, email, hashPassword, token, isActivated, activationLink, role) VALUE 
		(?, ?, ?, ?, ?, ?, ?)
	`
	var args []interface{}
	args = append(args,
		user.FullName, user.Email,
		user.HashPassword, user.Token,
		user.IsActivated, user.ActivationLink, user.Role,
	)

	result, err := ur.MySqlConn.Exec(sql, args...)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}

func (ur *UserRepository) SaveToken(userId int64, token string) error {
	sql := `
		UPDATE users
		SET token = ?
		WHERE id = ?
	`
	if _, err := ur.MySqlConn.Exec(sql, token, userId); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) SetPassword(userID int64, hashPassword string) error {
	sql := `
		UPDATE users
		SET hashPassword = ?
		WHERE id = ?
	`

	_, err := ur.MySqlConn.Exec(sql, hashPassword, userID)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) RemoveToken(token string) error {
	sql := `
		UPDATE users
		SET token = null
		WHERE token = ?
	`
	if _, err := ur.MySqlConn.Exec(sql, token); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) Activate(userId int64) error {
	sql := `
		UPDATE users
		SET isActivated = 1
		WHERE id = ?
	`
	if _, err := ur.MySqlConn.Exec(sql, userId); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) DeleteUser(email string) (int64, error) {
	sql := `
		DELETE FROM users
		WHERE email = ?
	`

	var count int64

	result, err := ur.MySqlConn.Exec(sql, email)
	if err != nil {
		return count, err
	}

	count, err = result.RowsAffected()
	if err != nil {
		return count, err
	}

	return count, nil
}
