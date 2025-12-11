package jwt

import "github.com/golang-jwt/jwt/v5"

type TokenType string

const (
	AccessType  TokenType = "access"
	RefreshType TokenType = "refresh"
)

type claims struct {
	UserID int64     `json:"user_id"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}
