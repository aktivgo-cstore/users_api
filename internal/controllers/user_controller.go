package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"users_api/internal/dto"
	"users_api/internal/helpers"
	"users_api/internal/service"
)

var (
	clientURL = os.Getenv("CLIENT_URL")
)

type UserController struct {
	UserService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) Registration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("unable to read request body: " + err.Error())
		helpers.ErrorResponse(w, "Некорректный запрос", http.StatusInternalServerError)
		return
	}

	var userRegistrationData *dto.UserRegistrationData
	if err = json.Unmarshal(body, &userRegistrationData); err != nil {
		log.Println("unable to decode request body: " + err.Error())
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	token, er := uc.UserService.Register(userRegistrationData)
	if er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, token)))
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("unable to read request body: " + err.Error())
		helpers.ErrorResponse(w, "Некорректный запрос", http.StatusInternalServerError)
		return
	}

	var userLoginData *dto.UserLoginData
	if err = json.Unmarshal(body, &userLoginData); err != nil {
		log.Println("unable to decode request body: " + err.Error())
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	token, userData, er := uc.UserService.Login(userLoginData.Email, userLoginData.Password)
	if er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(
		`{"token": "%s", "id": "%d", "email": "%s", "is_activated": "%d", "role": "%s"}`,
		token, userData.ID, userData.Email, userData.IsActivated, userData.Role)))
}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("unable to read request body: " + err.Error())
		helpers.ErrorResponse(w, "Некорректный запрос", http.StatusInternalServerError)
		return
	}

	var userLogoutData *dto.UserLogoutData
	if err = json.Unmarshal(body, &userLogoutData); err != nil {
		log.Println("unable to decode request body: " + err.Error())
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if er := uc.UserService.Logout(userLogoutData.Token); er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, er := service.GetToken(r.Header)
	if er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	if er = service.CheckAccess(token); er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	users, er := uc.UserService.GetUsers()
	if er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	encode, err := json.Marshal(users)
	if err != nil {
		log.Println("unable to encode users: " + err.Error())
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(encode)
}
