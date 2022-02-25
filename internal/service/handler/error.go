package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) ErrorResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, message)))
}
