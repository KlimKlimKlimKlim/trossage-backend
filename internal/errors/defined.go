package errors

import (
	"errors"
	"net/http"
)

var (
	ErrInternalServerError = New("Internal server error", http.StatusInternalServerError)

	ErrInvalidBody        = New("Invalid body", http.StatusBadRequest)
	ErrEmptyBody          = New("Empty body", http.StatusUnprocessableEntity)
	ErrEmptyQuery         = New("Empty query", http.StatusUnprocessableEntity)
	ErrInvalidLogin       = New("Login must be 3-20 characters long and contain only lowercase letters, digits, and underscores", http.StatusUnprocessableEntity)
	ErrInvalidPassword    = New("Password must be 8-63 characters long, contain at least one letter and one digit, and use only allowed characters", http.StatusUnprocessableEntity)
	ErrInvalidDisplayName = New("Display name must be 1-20 characters long and contain only letters, digits, spaces, dots, underscores, and hyphens", http.StatusUnprocessableEntity)

	ErrUnauthorized = New("Unauthorized", http.StatusUnauthorized)

	ErrUserNotFound      = New("User not found", http.StatusNotFound)
	ErrUserAlreadyExists = New("User already exists", http.StatusConflict)

	ErrTokenNotFound = New("Token not found", http.StatusNotFound)

	ErrSamePassword = New("New password must be different from current", http.StatusUnprocessableEntity)
)

var (
	ErrUserIDIsEmpty  = errors.New("user_id is empty")
	ErrTokenIDIsEmpty = errors.New("token_id is empty")
)
