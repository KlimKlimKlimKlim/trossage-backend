package jwt

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
)

type Controller struct {
	secret    []byte
	lifetime  time.Duration
	tokenType TokenType
}

func New(cfg *config.JWTSetting, tokenType TokenType) *Controller {
	return &Controller{
		secret:    []byte(cfg.Secret),
		lifetime:  cfg.Lifetime,
		tokenType: tokenType,
	}
}

func (jc *Controller) generateToken(userID int64) *jwt.Token {
	c := claims{
		UserID: userID,
		Type:   jc.tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jc.lifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, c)
}

func (jc *Controller) signToken(token *jwt.Token) (string, error) {
	tokenString, err := token.SignedString(jc.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (jc *Controller) GenerateSignedToken(userID int64) (string, error) {
	refreshToken := jc.generateToken(userID)

	return jc.signToken(refreshToken)
}

func (jc *Controller) GenerateSignedTokenAndModel(userID int64) (string, models.Token, error) {
	refreshToken := jc.generateToken(userID)

	expiredAt, err := refreshToken.Claims.GetExpirationTime()
	if err != nil {
		return "", models.Token{}, fmt.Errorf("failed to get expiration time: %w", err)
	}

	refreshString, err := jc.signToken(refreshToken)
	if err != nil {
		return "", models.Token{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	hash := sha256.Sum256([]byte(refreshString))
	tokenHash := hex.EncodeToString(hash[:])

	modelToken := models.Token{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiredAt.Time,
	}

	return refreshString, modelToken, nil
}

func (jc *Controller) ProcessToken(tokenString string) (int64, error) {
	c := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, c, func(token *jwt.Token) (interface{}, error) {
		return jc.secret, nil
	})

	if err != nil || !token.Valid || c.UserID == 0 || c.Type != jc.tokenType {
		return 0, derrors.ErrInvalidToken
	}

	return c.UserID, nil
}
