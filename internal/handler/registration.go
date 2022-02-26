package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"users_api/internal/models"
	"users_api/internal/service/helpers"
)

func (h *Handler) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.ErrorResponse(w, "unable to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	var user []*models.User
	if err := json.Unmarshal(body, &user); err != nil {
		helpers.ErrorResponse(w, "unable to decode request body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	candidate, err := h.UserRepository.GetUser(user[0].Email)
	if err != nil {
		helpers.ErrorResponse(w, "unable to check candidate: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if candidate != nil {
		helpers.ErrorResponse(w, "user with this email already exists", http.StatusInternalServerError)
		return
	}

	user[0].RefreshToken = "123"

	if err := h.UserRepository.SaveUser(user[0]); err != nil {
		helpers.ErrorResponse(w, "unable to save user: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
