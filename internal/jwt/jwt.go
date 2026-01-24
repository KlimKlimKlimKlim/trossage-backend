package jwt

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/model"
)

type Controller struct {
	tokenType TokenType
	secret    []byte
	lifetime  time.Duration
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
	token := jc.generateToken(userID)

	return jc.signToken(token)
}

func (jc *Controller) GenerateSignedTokenAndModel(userID int64) (string, model.Token, error) {
	token := jc.generateToken(userID)

	expiredAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		return "", model.Token{}, fmt.Errorf("failed to get expiration time: %w", err)
	}

	tokenString, err := jc.signToken(token)
	if err != nil {
		return "", model.Token{}, fmt.Errorf("failed to sign token: %w", err)
	}

	hash := sha256.Sum256([]byte(tokenString))
	tokenHash := hex.EncodeToString(hash[:])

	modelToken := model.Token{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiredAt.Time,
	}

	return tokenString, modelToken, nil
}

func (jc *Controller) ProcessToken(tokenString string) (int64, error) {
	c := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, c, func(token *jwt.Token) (any, error) {
		return jc.secret, nil
	})

	if err != nil || !token.Valid || c.UserID <= 0 || c.Type != jc.tokenType {
		return 0, derrors.ErrUnauthorized
	}

	return c.UserID, nil
}
