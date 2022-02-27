package service

import (
	"fmt"
	"users_api/internal/errors"
	"users_api/internal/models"
	"users_api/internal/repository"
)

type UserService struct {
	UserRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

func (us *UserService) Register(user *models.User) (string, *errors.ApiError) {
	candidate, err := us.UserRepository.GetUser(user.Email)
	if err != nil {
		return "", errors.InternalServerError(err)
	}

	if candidate != nil {
		return "", errors.BadRequestError(fmt.Sprintf("user with email %s user.Email already exists", user.Email), nil)
	}

	user.RefreshToken = "123"

	if err = us.UserRepository.SaveUser(user); err != nil {
		return "", errors.InternalServerError(err)
	}

	return user.RefreshToken, nil
}

func (us *UserService) GetUsers() ([]*models.User, *errors.ApiError) {
	users, err := us.UserRepository.GetUsers()
	if err != nil {
		return nil, errors.InternalServerError(err)
	}

	return users, nil
}

func (us *UserService) DeleteUser(id int) *errors.ApiError {
	count, err := us.UserRepository.DeleteUser(id)
	if err != nil {
		return errors.InternalServerError(err)
	}
	
	if count < 1 {
		return errors.BadRequestError(fmt.Sprintf("user with id %d is not exists", id), nil)
	}

	return nil
}
