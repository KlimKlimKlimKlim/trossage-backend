package errors

import (
	"net/http"
)

var (
	ErrInvalidBody = New("invalid body", http.StatusBadRequest)

	ErrInvalidToken       = New("invalid token", http.StatusForbidden)
	ErrInvalidCredentials = New("invalid credentials", http.StatusUnauthorized)

	ErrUserNotFound      = New("user not found", http.StatusNotFound)
	ErrUserAlreadyExists = New("user already exists", http.StatusConflict)
)
