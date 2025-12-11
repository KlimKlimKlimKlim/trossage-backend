package models

import (
	"time"
)

type User struct {
	ID           int64
	Login        string
	PasswordHash string
	DisplayName  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
