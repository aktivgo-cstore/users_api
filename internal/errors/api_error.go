package errors

import "net/http"

type ApiError struct {
	status  int
	message string
	error   error
}

func NewApiError(status int, message string, error error) *ApiError {
	return &ApiError{
		status:  status,
		message: message,
		error:   error,
	}
}

func UnauthorizedError() *ApiError {
	return NewApiError(http.StatusUnauthorized, "Пользователь не авторизован", nil)
}

func BadRequestError(message string, error error) *ApiError {
	return NewApiError(http.StatusBadRequest, message, error)
}

func InternalServerError(error error) *ApiError {
	return NewApiError(http.StatusInternalServerError, "Упс... Что-то пошло не так...", error)
}
