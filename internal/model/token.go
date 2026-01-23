package model

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

func (t *Token) IsRevoked() bool {
	return !t.RevokedAt.IsZero()
}

type JWTPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenRevocation struct {
	Reason  string
	Revoked bool
}
