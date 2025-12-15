package errors

import (
	"net/http"
)

var (
	ErrInvalidBody = New("invalid body", http.StatusBadRequest)

	ErrUnauthorized = New("unauthorized", http.StatusUnauthorized)

	ErrUserNotFound      = New("user not found", http.StatusNotFound)
	ErrUserAlreadyExists = New("user already exists", http.StatusConflict)

	ErrTokenNotFound = New("token not found", http.StatusNotFound)
)
