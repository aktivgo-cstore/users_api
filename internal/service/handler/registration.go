package handler

import (
	"encoding/json"
	"net/http"
	"users_api/internal/models"
)

func (h *Handler) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	/*body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.ErrorResponse(w, "unable to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}*/

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.ErrorResponse(w, "unable to decode request body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	/*if err = json.Unmarshal(, &user); err != nil {
		h.ErrorResponse(w, "unable to decode request body: "+err.Error(), http.StatusInternalServerError)
		return
	}*/

	if err := h.UserRepository.AddUser(user); err != nil {
		h.ErrorResponse(w, "unable to save user: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
