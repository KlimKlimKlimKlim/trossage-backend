package errors

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidBody        = New("invalid body", http.StatusBadRequest)
	ErrEmptyBody          = New("empty body", http.StatusUnprocessableEntity)
	ErrInvalidLogin       = New("login must be 3-20 characters long and contain only lowercase letters, digits, and underscores", http.StatusUnprocessableEntity)
	ErrInvalidPassword    = New("password must be 8-63 characters long, contain at least one letter and one digit, and use only allowed characters", http.StatusUnprocessableEntity)
	ErrInvalidDisplayName = New("display name must be 1-20 characters long and contain only letters, digits, spaces, dots, underscores, and hyphens", http.StatusUnprocessableEntity)

	ErrUnauthorized = New("unauthorized", http.StatusUnauthorized)

	ErrUserNotFound      = New("user not found", http.StatusNotFound)
	ErrUserAlreadyExists = New("user already exists", http.StatusConflict)

	ErrTokenNotFound = New("token not found", http.StatusNotFound)

	ErrSamePassword = New("new password must be different from current", http.StatusUnprocessableEntity)
)

var (
	ErrUserIDIsEmpty  = errors.New("user_id is empty")
	ErrTokenIDIsEmpty = errors.New("token_id is empty")
)
