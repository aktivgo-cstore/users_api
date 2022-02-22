package handler

import (
	"fmt"
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
	fmt.Println("Hello from registration")
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello from login")
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello from logout")
}

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello from delete")
}

func (h *Handler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello from refresh")
}

func (h *Handler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello from users")

}