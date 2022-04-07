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
	MailService    *MailService
	UserRepository *repository.UserRepository
}

func NewUserService(mailService *MailService, userRepository *repository.UserRepository) *UserService {
	return &UserService{
		MailService:    mailService,
		UserRepository: userRepository,
	}
}

func (us *UserService) Register(userData *dto.UserRegistrationData) (string, *errors.ApiError) {
	candidate, err := us.UserRepository.GetUserByEmail(userData.Email)
	if err != nil {
		log.Println("unable to get user by email: " + err.Error())
		return "", errors.InternalServerError(err)
	}

	if candidate != nil {
		return "", errors.BadRequestError(fmt.Sprintf("Пользователь с электронной почтой %s уже существует", userData.Email),
			fmt.Errorf("the user with the email %s already exists", userData.Email))
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), 3)
	if err != nil {
		log.Println("unable to hash password: " + err.Error())
		return "", errors.InternalServerError(err)
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
		return "", errors.InternalServerError(err)
	}

	tokenData := dto.NewTokenData(id, user.Email, user.IsActivated, user.Role)
	token, err := GenerateToken(tokenData)
	if err != nil {
		log.Println("unable to generate token: " + err.Error())
		return "", errors.InternalServerError(err)
	}

	if err = us.UserRepository.SaveToken(id, token); err != nil {
		log.Println("unable to save token: " + err.Error())
		return "", errors.InternalServerError(err)
	}

	if err = us.MailService.SendActivationMail(userData.Email, user.FullName, apiUrl+":"+apiPort+"/activate/"+activationLink); err != nil {
		log.Println("unable to send mail: " + err.Error())
		return "", errors.InternalServerError(err)
	}

	return token, nil
}

func (us *UserService) Activate(activationLink string) *errors.ApiError {
	user, err := us.UserRepository.GetUserByActivationLink(activationLink)
	if err != nil {
		log.Println("unable to get user by activation link: " + err.Error())
		return errors.InternalServerError(err)
	}

	if user == nil {
		return errors.BadRequestError("Некорректная ссылка активации", fmt.Errorf("invalid activation link"))
	}

	if err = us.UserRepository.Activate(user.ID); err != nil {
		log.Println("unable to activate user: " + err.Error())
		return errors.InternalServerError(err)
	}

	return nil
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

	tokenData := dto.NewTokenData(user.ID, user.Email, user.IsActivated, user.Role)

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

func (us *UserService) UpdateUser(userUpdateData *dto.UserUpdateData) *errors.ApiError {
	user, err := us.UserRepository.GetUserByID(userUpdateData.ID)
	if err != nil {
		log.Println("unable to get user by id: " + err.Error())
		return errors.InternalServerError(err)
	}

	if user == nil {
		return errors.BadRequestError(fmt.Sprintf("Пользователя с id %d не существует", userUpdateData.ID),
			fmt.Errorf("the user with the id %d does not exist", userUpdateData.ID))
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userUpdateData.Password), 3)
	if err != nil {
		log.Println("unable to hash password: " + err.Error())
		return errors.InternalServerError(err)
	}

	if err = us.UserRepository.SetPassword(user.ID, string(hashPassword)); err != nil {
		log.Println("unable to set temporary password: " + err.Error())
		return errors.InternalServerError(err)
	}

	return nil
}

func (us *UserService) DeleteUser(email string) *errors.ApiError {
	count, err := us.UserRepository.DeleteUser(email)
	if err != nil {
		log.Println("unable to delete user: " + err.Error())
		return errors.InternalServerError(err)
	}

	if count < 1 {
		return errors.BadRequestError(fmt.Sprintf("Пользователя с электронной почтой %s не существует", email),
			fmt.Errorf("the user with the email %s does not exist", email))
	}

	return nil
}

func (us *UserService) RestorePassword(userRestorePasswordData *dto.UserRestorePasswordData) *errors.ApiError {
	user, err := us.UserRepository.GetUserByEmail(userRestorePasswordData.Email)
	if err != nil {
		log.Println("unable to get user by email: " + err.Error())
		return errors.InternalServerError(err)
	}

	if user == nil {
		return errors.BadRequestError(fmt.Sprintf("Пользователя с электронной почтой %s не существует", userRestorePasswordData.Email),
			fmt.Errorf("the user with the email %s does not exist", userRestorePasswordData.Email))
	}

	tempPassword := uuid.NewString()

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), 3)
	if err != nil {
		log.Println("unable to hash password: " + err.Error())
		return errors.InternalServerError(err)
	}

	if err = us.UserRepository.SetPassword(user.ID, string(hashPassword)); err != nil {
		log.Println("unable to set temporary password: " + err.Error())
		return errors.InternalServerError(err)
	}

	if err = us.MailService.SendRestorePasswordMail(user.Email, user.FullName, tempPassword); err != nil {
		log.Println("unable to send mail: " + err.Error())
		return errors.InternalServerError(err)
	}

	return nil
}

func (us *UserService) Logout(token string) *errors.ApiError {
	if err := us.UserRepository.RemoveToken(token); err != nil {
		log.Println("unable to remove token: " + err.Error())
		return errors.InternalServerError(err)
	}

	return nil
}

func (us *UserService) GetUsers() ([]*models.User, *errors.ApiError) {
	users, err := us.UserRepository.GetUsers()
	if err != nil {
		log.Println("unable to get users: " + err.Error())
		return nil, errors.InternalServerError(err)
	}

	return users, nil
}
