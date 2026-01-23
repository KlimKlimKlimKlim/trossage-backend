package model

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

func (u *User) IsEmpty() bool {
	return u.ID == 0
}

type AuthUser struct {
	PasswordHash string
	User
}
