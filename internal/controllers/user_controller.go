package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
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

	var userData *dto.UserData
	if err = json.Unmarshal(body, &userData); err != nil {
		log.Println("unable to decode request body: " + err.Error())
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	tokens, er := uc.UserService.Register(userData)
	if er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"accessToken": "%s", "refreshToken": "%s"}`, tokens.AccessToken, tokens.RefreshToken)))
}

func (uc *UserController) Activate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	activationLink := mux.Vars(r)["link"]

	if er := uc.UserService.Activate(activationLink); er != nil {
		log.Println("unable to activate user: " + er.Error.Error())
		helpers.ErrorResponse(w, er.Message, er.Status)
	}

	http.Redirect(w, r, clientURL, http.StatusSeeOther)
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Hello from login")
	if err != nil {
		log.Fatalln("Unable to write string. Error: " + err.Error())
	}
}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Hello from logout")
	if err != nil {
		log.Fatalln("Unable to write string. Error: " + err.Error())
	}
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("unable to read request body: " + err.Error())
		helpers.ErrorResponse(w, "ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}

	var decoded map[string]string
	if err = json.Unmarshal(body, &decoded); err != nil {
		log.Println("unable to decode request body: " + err.Error())
		helpers.ErrorResponse(w, "ошибка при декодировании тела запроса", http.StatusInternalServerError)
		return
	}

	if er := uc.UserService.DeleteUser(decoded["email"]); er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (uc *UserController) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Hello from refresh")
	if err != nil {
		log.Fatalln("Unable to write string. Error: " + err.Error())
	}
}

func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
