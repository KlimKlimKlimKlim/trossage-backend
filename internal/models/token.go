package models

import (
	"time"
)

type Token struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
	RevokedAt time.Time
	CreatedAt time.Time
}
