package errors

import (
	"errors"
	"net/http"
)

var (
	ErrInternalServerError = New(http.StatusInternalServerError, "Internal server error")

	ErrInvalidBody  = New(http.StatusBadRequest, "Invalid body")
	ErrEmptyBody    = New(http.StatusUnprocessableEntity, "Empty body")
	ErrEmptyQuery   = New(http.StatusUnprocessableEntity, "Empty query")
	ErrInvalidLogin = New(
		http.StatusUnprocessableEntity,
		"Login must be 3-20 characters long and contain only lowercase letters, digits, and underscores",
	)
	ErrInvalidPassword = New(
		http.StatusUnprocessableEntity,
		"Password must be 8-63 characters long, contain at least one letter and one digit, and use only allowed characters",
	)
	ErrInvalidDisplayName = New(
		http.StatusUnprocessableEntity,
		"Display name must be 1-20 characters long and contain only letters, digits, spaces, dots, underscores, and hyphens",
	)
	ErrSearchQueryTooShort = New(http.StatusUnprocessableEntity, "Search query too short")
	ErrInvalidPathParams   = New(http.StatusBadRequest, "Invalid path params")

	ErrUnauthorized = New(http.StatusUnauthorized, "Unauthorized")

	ErrUserNotFound      = New(http.StatusNotFound, "User not found")
	ErrUserAlreadyExists = New(http.StatusConflict, "User already exists")

	ErrTokenNotFound = New(http.StatusNotFound, "Token not found")

	ErrSamePassword = New(http.StatusUnprocessableEntity, "New password must be different from current")

	ErrCannotChatWithYourself = New(http.StatusUnprocessableEntity, "Cannot create chat with yourself")
	ErrChatAlreadyExists      = New(http.StatusConflict, "Chat already exists")
	ErrChatNotFound           = New(http.StatusNotFound, "Chat not found")
	ErrUserIsNotMember        = New(http.StatusForbidden, "User is not member")

	ErrMessageIsEmpty = New(http.StatusUnprocessableEntity, "Message is empty")
)

var (
	ErrUserIDIsEmpty    = errors.New("user_id is empty")
	ErrTokenIDIsEmpty   = errors.New("token_id is empty")
	ErrUnknownLoggerEnv = errors.New("unknown logger environment")
	ErrPanicInTx        = errors.New("panic in transaction")
	ErrEmptyInput       = errors.New("empty input")
)
