package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"users_api/internal/repository"
)

type Handler struct {
	UserRepository *repository.UserRepository
}

func NewHandler(userRepository *repository.UserRepository) *Handler {
	return &Handler{
		UserRepository: userRepository,
	}
}

func (h *Handler) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Hello from registration")
	if err != nil {
		log.Fatalln("Unable to write string. Error: " + err.Error())
	}
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Hello from login")
	if err != nil {
		log.Fatalln("Unable to write string. Error: " + err.Error())
	}
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Hello from logout")
	if err != nil {
		log.Fatalln("Unable to write string. Error: " + err.Error())
	}
}

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Hello from delete")
	if err != nil {
		log.Fatalln("Unable to write string. Error: " + err.Error())
	}
}

func (h *Handler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Hello from refresh")
	if err != nil {
		log.Fatalln("Unable to write string. Error: " + err.Error())
	}
}

func (h *Handler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Hello from users")
	if err != nil {
		log.Fatalln("Unable to write string. Error: " + err.Error())
	}
}