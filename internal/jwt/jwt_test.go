package jwt

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"

	"github.com/GlaciemArgentum/trossage-backend/internal/config"
	derrors "github.com/GlaciemArgentum/trossage-backend/internal/errors"
)

type JWTTestSuite struct {
	suite.Suite
}

func TestJWTTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(JWTTestSuite))
}

func (s *JWTTestSuite) TestNew() {
	testCases := []struct {
		name             string
		config           *config.JWTSetting
		tokenType        TokenType
		expectedLifetime time.Duration
	}{
		{
			name: "success - creates access token controller",
			config: &config.JWTSetting{
				Secret:   "test_secret_access",
				Lifetime: 15 * time.Minute,
			},
			tokenType:        AccessType,
			expectedLifetime: 15 * time.Minute,
		},
		{
			name: "success - creates refresh token controller",
			config: &config.JWTSetting{
				Secret:   "test_secret_refresh",
				Lifetime: 7 * 24 * time.Hour,
			},
			tokenType:        RefreshType,
			expectedLifetime: 7 * 24 * time.Hour,
		},
		{
			name: "success - creates controller with short lifetime",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 1 * time.Second,
			},
			tokenType:        AccessType,
			expectedLifetime: 1 * time.Second,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.T().Parallel()

			// Test.
			controller := New(tc.config, tc.tokenType)

			// Check.
			s.Require().NotNil(controller)
			s.Equal([]byte(tc.config.Secret), controller.secret)
			s.Equal(tc.expectedLifetime, controller.lifetime)
			s.Equal(tc.tokenType, controller.tokenType)
		})
	}
}

func (s *JWTTestSuite) TestGenerateSignedToken() {
	testCases := []struct {
		name      string
		config    *config.JWTSetting
		tokenType TokenType
		userID    int64
	}{
		{
			name: "success - generates access token for valid user",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 15 * time.Minute,
			},
			tokenType: AccessType,
			userID:    12345,
		},
		{
			name: "success - generates refresh token for valid user",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 7 * 24 * time.Hour,
			},
			tokenType: RefreshType,
			userID:    67890,
		},
		{
			name: "success - generates token for max int64 user ID",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 1 * time.Hour,
			},
			tokenType: AccessType,
			userID:    9223372036854775807,
		},
		{
			name: "success - generates token for min positive user ID",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 1 * time.Hour,
			},
			tokenType: AccessType,
			userID:    1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.T().Parallel()

			controller := New(tc.config, tc.tokenType)

			// Test.
			tokenString, err := controller.GenerateSignedToken(tc.userID)

			// Check.
			s.Require().NoError(err)
			s.NotEmpty(tokenString)
			s.Contains(tokenString, ".")

			parts := strings.Split(tokenString, ".")
			s.Len(parts, 3)

			parsedToken, err := jwt.ParseWithClaims(tokenString, &claims{}, func(_ *jwt.Token) (any, error) {
				return controller.secret, nil
			})
			s.Require().NoError(err)
			s.True(parsedToken.Valid)

			parsedClaims, ok := parsedToken.Claims.(*claims)
			s.Require().True(ok)
			s.Equal(tc.userID, parsedClaims.UserID)
			s.Equal(tc.tokenType, parsedClaims.Type)
		})
	}
}

func (s *JWTTestSuite) TestGenerateSignedTokenAndModel() {
	testCases := []struct {
		name      string
		config    *config.JWTSetting
		tokenType TokenType
		userID    int64
	}{
		{
			name: "success - generates access token and model",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 15 * time.Minute,
			},
			tokenType: AccessType,
			userID:    12345,
		},
		{
			name: "success - generates refresh token and model",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 7 * 24 * time.Hour,
			},
			tokenType: RefreshType,
			userID:    67890,
		},
		{
			name: "success - token hash is consistent",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 1 * time.Hour,
			},
			tokenType: AccessType,
			userID:    999,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.T().Parallel()

			controller := New(tc.config, tc.tokenType)

			// Test.
			tokenString, modelToken, err := controller.GenerateSignedTokenAndModel(tc.userID)

			// Check.
			s.Require().NoError(err)
			s.NotEmpty(tokenString)
			s.Equal(tc.userID, modelToken.UserID)
			s.NotEmpty(modelToken.TokenHash)
			s.NotZero(modelToken.ExpiresAt)

			s.Len(modelToken.TokenHash, 64)

			expectedHash := sha256.Sum256([]byte(tokenString))
			expectedHashString := hex.EncodeToString(expectedHash[:])
			s.Equal(expectedHashString, modelToken.TokenHash)

			expectedExpiry := time.Now().Add(tc.config.Lifetime)
			s.WithinDuration(expectedExpiry, modelToken.ExpiresAt, 2*time.Second)

			parsedToken, err := jwt.ParseWithClaims(tokenString, &claims{}, func(_ *jwt.Token) (any, error) {
				return controller.secret, nil
			})
			s.Require().NoError(err)

			parsedClaims, ok := parsedToken.Claims.(*claims)
			s.Require().True(ok)
			s.Equal(tc.userID, parsedClaims.UserID)
			s.Equal(tc.tokenType, parsedClaims.Type)
		})
	}
}

func (s *JWTTestSuite) TestTokenLifetime() {
	testCases := []struct {
		config   *config.JWTSetting
		name     string
		lifetime time.Duration
	}{
		{
			name: "success - short lifetime token has correct expiry",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 1 * time.Minute,
			},
			lifetime: 1 * time.Minute,
		},
		{
			name: "success - long lifetime token has correct expiry",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 7 * 24 * time.Hour,
			},
			lifetime: 7 * 24 * time.Hour,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.T().Parallel()

			controller := New(tc.config, AccessType)
			beforeGeneration := time.Now()

			// Test.
			tokenString, err := controller.GenerateSignedToken(123)
			s.Require().NoError(err)

			afterGeneration := time.Now()

			// Parse and check claims.
			parsedToken, err := jwt.ParseWithClaims(tokenString, &claims{}, func(_ *jwt.Token) (any, error) {
				return controller.secret, nil
			})
			s.Require().NoError(err)

			parsedClaims, ok := parsedToken.Claims.(*claims)
			s.Require().True(ok)

			// Check IssuedAt is within generation time window.
			s.Require().NotNil(parsedClaims.IssuedAt)
			issuedAt := parsedClaims.IssuedAt.Time
			issuedAtValid := issuedAt.After(beforeGeneration.Add(-1*time.Second)) &&
				issuedAt.Before(afterGeneration.Add(1*time.Second))
			s.True(issuedAtValid, "IssuedAt should be within generation time window")

			// Check ExpiresAt is IssuedAt + lifetime.
			s.Require().NotNil(parsedClaims.ExpiresAt)
			expiresAt := parsedClaims.ExpiresAt.Time
			expectedExpiry := issuedAt.Add(tc.lifetime)
			s.WithinDuration(expectedExpiry, expiresAt, 1*time.Second,
				"ExpiresAt should be IssuedAt + lifetime")
		})
	}
}

func (s *JWTTestSuite) TestProcessToken() {
	testCases := []struct {
		expectedErr  error
		config       *config.JWTSetting
		modifyToken  func(string) string
		name         string
		tokenType    TokenType
		userID       int64
		expectUserID int64
		wrongSecret  bool
		wrongType    bool
	}{
		{
			name: "success - processes valid access token",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 15 * time.Minute,
			},
			tokenType:    AccessType,
			userID:       12345,
			expectUserID: 12345,
		},
		{
			name: "success - processes valid refresh token",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 7 * 24 * time.Hour,
			},
			tokenType:    RefreshType,
			userID:       67890,
			expectUserID: 67890,
		},
		{
			name: "error - token with wrong secret",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 1 * time.Hour,
			},
			tokenType:   AccessType,
			userID:      123,
			wrongSecret: true,
			expectedErr: derrors.ErrUnauthorized,
		},
		{
			name: "error - expired token",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: -1 * time.Hour,
			},
			tokenType:   AccessType,
			userID:      123,
			expectedErr: derrors.ErrUnauthorized,
		},
		{
			name: "error - wrong token type",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 1 * time.Hour,
			},
			tokenType:   AccessType,
			userID:      123,
			wrongType:   true,
			expectedErr: derrors.ErrUnauthorized,
		},
		{
			name: "error - user ID is zero",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 1 * time.Hour,
			},
			tokenType:   AccessType,
			userID:      0,
			expectedErr: derrors.ErrUnauthorized,
		},
		{
			name: "error - negative user ID",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 1 * time.Hour,
			},
			tokenType: AccessType,
			userID:    -1,
			modifyToken: func(_ string) string {
				c := claims{
					UserID: -1,
					Type:   AccessType,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
				tokenString, _ := token.SignedString([]byte("test_secret"))

				return tokenString
			},
			expectedErr: derrors.ErrUnauthorized,
		},
		{
			name: "error - malformed token",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 1 * time.Hour,
			},
			tokenType: AccessType,
			userID:    123,
			modifyToken: func(_ string) string {
				return "not.a.valid.token"
			},
			expectedErr: derrors.ErrUnauthorized,
		},
		{
			name: "error - empty token",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 1 * time.Hour,
			},
			tokenType: AccessType,
			userID:    123,
			modifyToken: func(_ string) string {
				return ""
			},
			expectedErr: derrors.ErrUnauthorized,
		},
		{
			name: "error - token with modified signature",
			config: &config.JWTSetting{
				Secret:   "test_secret",
				Lifetime: 1 * time.Hour,
			},
			tokenType: AccessType,
			userID:    123,
			modifyToken: func(token string) string {
				parts := strings.Split(token, ".")
				if len(parts) == 3 {
					parts[2] = "invalidsignature"
					return strings.Join(parts, ".")
				}

				return token
			},
			expectedErr: derrors.ErrUnauthorized,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.T().Parallel()

			controller := New(tc.config, tc.tokenType)

			var (
				tokenString string
				err         error
			)

			if tc.modifyToken == nil {
				generateTokenType := tc.tokenType
				if tc.wrongType {
					if tc.tokenType == AccessType {
						generateTokenType = RefreshType
					} else {
						generateTokenType = AccessType
					}
				}

				generateController := controller

				if tc.wrongSecret {
					wrongConfig := &config.JWTSetting{
						Secret:   "wrong_secret",
						Lifetime: tc.config.Lifetime,
					}
					generateController = New(wrongConfig, generateTokenType)
				} else if tc.wrongType {
					generateController = New(tc.config, generateTokenType)
				}

				tokenString, err = generateController.GenerateSignedToken(tc.userID)
				s.Require().NoError(err)
			} else {
				tokenString, err = controller.GenerateSignedToken(tc.userID)
				s.Require().NoError(err)

				tokenString = tc.modifyToken(tokenString)
			}

			// Test.
			userID, err := controller.ProcessToken(tokenString)

			// Check.
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
				s.Zero(userID)

				return
			}

			s.Require().NoError(err)
			s.Equal(tc.expectUserID, userID)
		})
	}
}
