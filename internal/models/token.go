package models

import (
	"time"
)

type Token struct {
	ExpiresAt time.Time
	RevokedAt time.Time
	CreatedAt time.Time
	TokenHash string
	ID        int64
	UserID    int64
}

type JWTPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenRevocation struct {
	Reason  string
	Revoked bool
}
