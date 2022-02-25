package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	users, err := h.UserRepository.GetUsers()
	if err != nil {
		h.ErrorResponse(w, "unable to get users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(users); err != nil {
		h.ErrorResponse(w, "unable to encode users: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
