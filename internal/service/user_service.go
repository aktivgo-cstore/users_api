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

func (us *UserService) Register(userData *dto.UserRegistrationData) (int, string, *errors.ApiError) {
	candidate, err := us.UserRepository.GetUserByEmail(userData.Email)
	if err != nil {
		log.Println("unable to get user by email: " + err.Error())
		return -1, "", errors.InternalServerError(err)
	}

	if candidate != nil {
		return -1, "", errors.BadRequestError(fmt.Sprintf("Пользователь с электронной почтой %s уже существует", userData.Email),
			fmt.Errorf("the user with the email %s already exists", userData.Email))
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), 3)
	if err != nil {
		log.Println("unable to hash password: " + err.Error())
		return -1, "", errors.InternalServerError(err)
	}

	activationLink := uuid.New().String()

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
		log.Println("unable to save user: " + err.Error())
		return -1, "", errors.InternalServerError(err)
	}

	tokenData := dto.NewTokenData(int(id), user.Email, user.IsActivated, user.Role)
	token, err := GenerateToken(tokenData)
	if err != nil {
		log.Println("unable to generate token: " + err.Error())
		return -1, "", errors.InternalServerError(err)
	}

	if err = us.UserRepository.SaveToken(id, token); err != nil {
		log.Println("unable to save token: " + err.Error())
		return -1, "", errors.InternalServerError(err)
	}

	return int(id), token, nil
}

func (us *UserService) Login(email string, password string) (string, *dto.TokenData, *errors.ApiError) {
	user, err := us.UserRepository.GetUserByEmail(email)
	if err != nil {
		log.Println("unable to get user by email: " + err.Error())
		return "", nil, errors.InternalServerError(err)
	}

	if user == nil {
		return "", nil, errors.BadRequestError("Пользователь с таким email не найден",
			fmt.Errorf("the user with this email was not found"))
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(password)); err != nil {
		return "", nil, errors.BadRequestError("Неверный пароль",
			fmt.Errorf("invalid password"))
	}

	tokenData := dto.NewTokenData(int(user.ID), user.Email, user.IsActivated, user.Role)

	token, err := GenerateToken(tokenData)
	if err != nil {
		log.Println("unable to generate token: " + err.Error())
		return "", nil, errors.InternalServerError(err)
	}

	if err = us.UserRepository.SaveToken(user.ID, token); err != nil {
		log.Println("unable to save token: " + err.Error())
		return "", nil, errors.InternalServerError(err)
	}

	return token, tokenData, nil
}

func (us *UserService) Logout(token string) *errors.ApiError {
	if err := us.UserRepository.RemoveToken(token); err != nil {
		log.Println("unable to remove token: " + err.Error())
		return errors.InternalServerError(err)
	}

	return nil
}

func (us *UserService) GetUserByID(userID int) (*models.User, *errors.ApiError) {
	user, err := us.UserRepository.GetUserByID(int64(userID))
	if err != nil {
		log.Println("unable to get user: " + err.Error())
		return nil, errors.InternalServerError(err)
	}

	return user, nil
}

func (us *UserService) GetUsers() ([]*models.User, *errors.ApiError) {
	users, err := us.UserRepository.GetUsers()
	if err != nil {
		log.Println("unable to get users: " + err.Error())
		return nil, errors.InternalServerError(err)
	}

	return users, nil
}
