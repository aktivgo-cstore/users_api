package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"users_api/internal/service/helpers"
)

type DeletedUser struct {
	ID int `json:"id"`
}

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.ErrorResponse(w, "unable to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	var deletedUser []*DeletedUser
	if err := json.Unmarshal(body, &deletedUser); err != nil {
		helpers.ErrorResponse(w, "unable to decode request body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.UserRepository.DeleteUser(deletedUser[0].ID); err != nil {
		helpers.ErrorResponse(w, "unable to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
