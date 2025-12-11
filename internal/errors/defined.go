package errors

import (
	"net/http"
)

var (
	ErrInvalidBody = New("invalid body", http.StatusBadRequest)

	ErrInvalidToken         = New("invalid token", http.StatusForbidden)
	ErrWrongLoginOrPassword = New("wrong login or password", http.StatusForbidden)

	ErrUserNotFound      = New("user not found", http.StatusNotFound)
	ErrUserAlreadyExists = New("user already exists", http.StatusConflict)
)
