package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"users_api/internal/helpers"
	"users_api/internal/models"
	"users_api/internal/service"
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
	w.WriteHeader(http.StatusCreated)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("unable to read request body: " + err.Error())
		helpers.ErrorResponse(w, "unable to read request body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var user *models.User
	if err := json.Unmarshal(body, &user); err != nil {
		helpers.ErrorResponse(w, "unable to decode request body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, error := uc.UserService.Register(user)
	if err != nil {
		helpers.ErrorResponse(w, "unable to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte(fmt.Sprintf(`{"accessToken": "%s"}`, accessToken)))
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
	w.WriteHeader(http.StatusOK)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.ErrorResponse(w, "unable to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	var decoded map[string]int
	if err := json.Unmarshal(body, &decoded); err != nil {
		helpers.ErrorResponse(w, "unable to decode request body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := uc.UserService.DeleteUser(decoded["id"]); err != nil {
		helpers.ErrorResponse(w, "unable to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (uc *UserController) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Hello from refresh")
	if err != nil {
		log.Fatalln("Unable to write string. Error: " + err.Error())
	}
}

func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	users, err := uc.UserService.GetUsers()
	if err != nil {
		helpers.ErrorResponse(w, "unable to get users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(users); err != nil {
		helpers.ErrorResponse(w, "unable to encode users: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
