package models

import (
	"time"
)

type User struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	Login       string
	DisplayName string
	ID          int64
}

func (u *User) IsDeleted() bool {
	return !u.DeletedAt.IsZero()
}

type AuthUser struct {
	PasswordHash string
	User
}
