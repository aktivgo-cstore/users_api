package service

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"users_api/internal/dto"
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

func (us *UserService) Register(userData *dto.UserData) (string, *errors.ApiError) {
	candidate, err := us.UserRepository.GetUserByEmail(userData.Email)
	if err != nil {
		return "", errors.InternalServerError(err)
	}

	if candidate != nil {
		return "", errors.BadRequestError(fmt.Sprintf("Пользователь с электронной почтой %s уже существует", userData.Email),
			fmt.Errorf("пользователь с электронной почтой %s уже существует", userData.Email))
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), 3)
	if err != nil {
		return "", errors.InternalServerError(err)
	}

	activationLink := uuid.New().String()

	if err = SendActivationMail(userData.Email, apiUrl+":"+apiPort+"/activate/"+activationLink); err != nil {
		log.Println("unable to send mail: " + err.Error())
		return "", errors.InternalServerError(err)
	}

	user := &models.User{
		FullName:       userData.FullName,
		Email:          userData.Email,
		HashPassword:   string(hashPassword),
		IsActivated:    0,
		ActivationLink: activationLink,
		Role:           "user",
	}

	id, err := us.UserRepository.SaveUser(user)
	if err != nil {
		return "", errors.InternalServerError(err)
	}

	tokenData := dto.NewTokenData(id, user.Email, user.Role)
	token, err := GenerateToken(tokenData)

	if err = us.UserRepository.SaveToken(id, token); err != nil {
		return "", errors.InternalServerError(err)
	}

	return token, nil
}

func (us *UserService) Activate(activationLink string) *errors.ApiError {
	user, err := us.UserRepository.GetUserByActivationLink(activationLink)
	if err != nil {
		return errors.InternalServerError(err)
	}

	if user == nil {
		return errors.BadRequestError("Неккоректная ссылка активации", fmt.Errorf("неккоректная ссылка активации"))
	}

	if err = us.UserRepository.Activate(user.ID); err != nil {
		return errors.InternalServerError(err)
	}

	return nil
}

func (us *UserService) GetUsers() ([]*models.User, *errors.ApiError) {
	users, err := us.UserRepository.GetUsers()
	if err != nil {
		return nil, errors.InternalServerError(err)
	}

	return users, nil
}

func (us *UserService) DeleteUser(email string) *errors.ApiError {
	count, err := us.UserRepository.DeleteUser(email)
	if err != nil {
		return errors.InternalServerError(err)
	}

	if count < 1 {
		return errors.BadRequestError(fmt.Sprintf("Пользователя с электронной почтой %s не существует", email),
			fmt.Errorf("пользователя с электронной почтой %s не существует", email))
	}

	return nil
}
