package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgerrcode"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/model"
)

func (s *RepositoryTestSuite) TestInsertRefreshToken() {
	testCases := []struct {
		name        string
		loginSuffix string
		pgCode      string
		token       model.Token
		setupUser   bool
		duplicate   bool
	}{
		{
			name:        "success - inserts valid refresh token",
			loginSuffix: "token1",
			setupUser:   true,
			token: model.Token{
				TokenHash: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
		},
		{
			name:        "success - multiple tokens for same user",
			loginSuffix: "token2",
			setupUser:   true,
			token: model.Token{
				TokenHash: "b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3",
				ExpiresAt: time.Now().Add(48 * time.Hour),
			},
		},
		{
			name:        "error - duplicate token hash",
			loginSuffix: "token3",
			setupUser:   true,
			duplicate:   true,
			token: model.Token{
				TokenHash: "d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5",
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
			pgCode: pgerrcode.UniqueViolation,
		},
		{
			name:        "error - invalid token hash length too short",
			loginSuffix: "token4",
			setupUser:   true,
			token: model.Token{
				TokenHash: "short",
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
			pgCode: pgerrcode.CheckViolation,
		},
		{
			name:        "error - invalid token hash length too long",
			loginSuffix: "token5",
			setupUser:   true,
			token: model.Token{
				TokenHash: "toolong12345678901234567890123456789012345678901234567890123456789012345",
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
			pgCode: pgerrcode.CheckViolation,
		},
		{
			name:        "error - user does not exist",
			loginSuffix: "token6",
			setupUser:   false,
			token: model.Token{
				UserID:    999999,
				TokenHash: "c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4",
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
			pgCode: pgerrcode.ForeignKeyViolation,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			token := tc.token
			if tc.setupUser {
				insertedUser := s.createTestUser(ctx, "tokenuser_"+tc.loginSuffix)
				token.UserID = insertedUser.ID

				if tc.duplicate {
					dupToken := model.Token{
						UserID:    insertedUser.ID,
						TokenHash: token.TokenHash,
						ExpiresAt: time.Now().Add(12 * time.Hour),
					}
					_, err := s.repo.InsertRefreshToken(ctx, dupToken)
					s.Require().NoError(err)
				}
			}

			// Test.
			result, err := s.repo.InsertRefreshToken(ctx, token)

			// Check.
			if tc.pgCode != "" {
				s.assertPgError(err, tc.pgCode)
				s.Zero(result.ID)

				return
			}

			s.Require().NoError(err)
			s.NotZero(result.ID)
			s.Equal(token.UserID, result.UserID)
			s.Equal(token.TokenHash, result.TokenHash)
			s.NotZero(result.CreatedAt)
			s.WithinDuration(token.ExpiresAt, result.ExpiresAt, time.Second)
			s.True(result.RevokedAt.IsZero())
		})
	}
}

func (s *RepositoryTestSuite) TestSelectRefreshToken() {
	testCases := []struct {
		expectedErr   error
		name          string
		loginSuffix   string
		tokenHash     string
		wrongUserID   bool
		wrongHash     bool
		revokeFirst   bool
		expectRevoked bool
	}{
		{
			name:          "success - finds valid refresh token",
			loginSuffix:   "select1",
			tokenHash:     "1111111111111111111111111111111111111111111111111111111111111111",
			wrongUserID:   false,
			wrongHash:     false,
			revokeFirst:   false,
			expectRevoked: false,
			expectedErr:   nil,
		},
		{
			name:          "error - token not found with wrong user id",
			loginSuffix:   "select2",
			tokenHash:     "2222222222222222222222222222222222222222222222222222222222222222",
			wrongUserID:   true,
			wrongHash:     false,
			revokeFirst:   false,
			expectRevoked: false,
			expectedErr:   derrors.ErrTokenNotFound,
		},
		{
			name:          "error - token not found with wrong token hash",
			loginSuffix:   "select3",
			tokenHash:     "3333333333333333333333333333333333333333333333333333333333333333",
			wrongUserID:   false,
			wrongHash:     true,
			revokeFirst:   false,
			expectRevoked: false,
			expectedErr:   derrors.ErrTokenNotFound,
		},
		{
			name:          "success - finds revoked token",
			loginSuffix:   "select4",
			tokenHash:     "4444444444444444444444444444444444444444444444444444444444444444",
			wrongUserID:   false,
			wrongHash:     false,
			revokeFirst:   true,
			expectRevoked: true,
			expectedErr:   nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			user := s.createTestUser(ctx, "selecttoken_"+tc.loginSuffix)
			insertedToken := s.createTestToken(ctx, user.ID, tc.tokenHash, time.Now().Add(24*time.Hour))

			if tc.revokeFirst {
				err := s.repo.RevokeRefreshTokenByID(ctx, insertedToken.ID)
				s.Require().NoError(err)
			}

			queryUserID := user.ID
			if tc.wrongUserID {
				queryUserID = 999999
			}

			queryHash := tc.tokenHash
			if tc.wrongHash {
				queryHash = "9999999999999999999999999999999999999999999999999999999999999999"
			}

			// Test.
			result, err := s.repo.SelectRefreshToken(ctx, queryUserID, queryHash)

			// Check.
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
				s.Zero(result.ID)

				return
			}

			s.Require().NoError(err)
			s.Equal(insertedToken.ID, result.ID)
			s.Equal(insertedToken.UserID, result.UserID)
			s.Equal(insertedToken.TokenHash, result.TokenHash)
			s.Equal(insertedToken.CreatedAt, result.CreatedAt)
			s.WithinDuration(insertedToken.ExpiresAt, result.ExpiresAt, time.Second)

			if tc.expectRevoked {
				s.False(result.RevokedAt.IsZero())
			} else {
				s.True(result.RevokedAt.IsZero())
			}
		})
	}
}

func (s *RepositoryTestSuite) TestRevokeRefreshTokenByID() {
	testCases := []struct {
		expectedErr   error
		name          string
		loginSuffix   string
		tokenHash     string
		alreadyRevoke bool
		invalidID     bool
	}{
		{
			name:          "success - revokes valid token",
			loginSuffix:   "revoke1",
			tokenHash:     "abc1111111111111111111111111111111111111111111111111111111111111",
			alreadyRevoke: false,
			invalidID:     false,
			expectedErr:   nil,
		},
		{
			name:          "error - token already revoked",
			loginSuffix:   "revoke2",
			tokenHash:     "abc2222222222222222222222222222222222222222222222222222222222222",
			alreadyRevoke: true,
			invalidID:     false,
			expectedErr:   derrors.ErrTokenNotFound,
		},
		{
			name:          "error - token not found",
			loginSuffix:   "revoke3",
			tokenHash:     "abc3333333333333333333333333333333333333333333333333333333333333",
			alreadyRevoke: false,
			invalidID:     true,
			expectedErr:   derrors.ErrTokenNotFound,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			user := s.createTestUser(ctx, "revoketoken_"+tc.loginSuffix)
			insertedToken := s.createTestToken(ctx, user.ID, tc.tokenHash, time.Now().Add(24*time.Hour))

			if tc.alreadyRevoke {
				err := s.repo.RevokeRefreshTokenByID(ctx, insertedToken.ID)
				s.Require().NoError(err)
			}

			revokeID := insertedToken.ID
			if tc.invalidID {
				revokeID = 999999
			}

			// Test.
			err := s.repo.RevokeRefreshTokenByID(ctx, revokeID)

			// Check.
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
				return
			}

			s.Require().NoError(err)

			revokedToken, err := s.repo.SelectRefreshToken(ctx, user.ID, tc.tokenHash)
			s.Require().NoError(err)
			s.False(revokedToken.RevokedAt.IsZero())
			s.True(revokedToken.RevokedAt.After(insertedToken.CreatedAt))
		})
	}
}

func (s *RepositoryTestSuite) TestRevokeRefreshTokensByUserID() {
	testCases := []struct {
		name            string
		loginSuffix     string
		createTokens    int
		revokeFirst     bool
		expectAllRevoke bool
	}{
		{
			name:            "success - revokes all user tokens",
			loginSuffix:     "revall1",
			createTokens:    3,
			revokeFirst:     false,
			expectAllRevoke: true,
		},
		{
			name:            "success - revokes only active tokens",
			loginSuffix:     "revall2",
			createTokens:    3,
			revokeFirst:     true,
			expectAllRevoke: true,
		},
		{
			name:            "success - no tokens to revoke",
			loginSuffix:     "revall3",
			createTokens:    0,
			revokeFirst:     false,
			expectAllRevoke: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			user := s.createTestUser(ctx, tc.loginSuffix)

			tokens := make([]model.Token, 0, tc.createTokens)
			for i := range tc.createTokens {
				tokenHash := s.generateTokenHash(tc.loginSuffix, i)
				insertedToken := s.createTestToken(ctx, user.ID, tokenHash, time.Now().Add(24*time.Hour))
				tokens = append(tokens, insertedToken)
			}

			if tc.revokeFirst && len(tokens) > 0 {
				err := s.repo.RevokeRefreshTokenByID(ctx, tokens[0].ID)
				s.Require().NoError(err)
			}

			// Test.
			err := s.repo.RevokeRefreshTokensByUserID(ctx, user.ID)

			// Check.
			s.Require().NoError(err)

			if tc.expectAllRevoke {
				for _, token := range tokens {
					result, err := s.repo.SelectRefreshToken(ctx, user.ID, token.TokenHash)
					s.Require().NoError(err)
					s.False(result.RevokedAt.IsZero())
				}
			}
		})
	}
}

func (s *RepositoryTestSuite) TestDeleteExpiredTokens() {
	testCases := []struct {
		name           string
		loginSuffix    string
		setupTokens    []time.Duration
		revokeFirst    int
		olderThan      time.Duration
		expectedDelete int64
	}{
		{
			name:           "success - deletes expired tokens",
			loginSuffix:    "expire1",
			setupTokens:    []time.Duration{-48 * time.Hour, -24 * time.Hour, 24 * time.Hour},
			revokeFirst:    -1,
			olderThan:      0,
			expectedDelete: 2,
		},
		{
			name:           "success - does not delete active tokens",
			loginSuffix:    "expire2",
			setupTokens:    []time.Duration{24 * time.Hour, 48 * time.Hour},
			revokeFirst:    -1,
			olderThan:      0,
			expectedDelete: 0,
		},
		{
			name:           "success - does not delete revoked tokens",
			loginSuffix:    "expire3",
			setupTokens:    []time.Duration{-48 * time.Hour, -24 * time.Hour},
			revokeFirst:    0,
			olderThan:      0,
			expectedDelete: 1,
		},
		{
			name:           "success - no tokens to delete",
			loginSuffix:    "expire4",
			setupTokens:    []time.Duration{},
			revokeFirst:    -1,
			olderThan:      0,
			expectedDelete: 0,
		},
		{
			name:           "success - custom threshold",
			loginSuffix:    "expire5",
			setupTokens:    []time.Duration{-72 * time.Hour, -48 * time.Hour, -24 * time.Hour},
			revokeFirst:    -1,
			olderThan:      -50 * time.Hour,
			expectedDelete: 1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			user := s.createTestUser(ctx, "expireuser_"+tc.loginSuffix)

			for i, expireOffset := range tc.setupTokens {
				tokenHash := s.generateTokenHash(tc.loginSuffix, i)
				insertedToken := s.createTestToken(ctx, user.ID, tokenHash, time.Now().Add(expireOffset))

				if tc.revokeFirst == i {
					err := s.repo.RevokeRefreshTokenByID(ctx, insertedToken.ID)
					s.Require().NoError(err)
				}
			}

			// Test.
			deletedCount, err := s.repo.DeleteExpiredTokens(ctx, time.Now().Add(tc.olderThan))

			// Check.
			s.Require().NoError(err)
			s.Equal(tc.expectedDelete, deletedCount)
		})
	}
}

func (s *RepositoryTestSuite) TestDeleteRevokedTokens() {
	testCases := []struct {
		name           string
		loginSuffix    string
		revokeIndexes  []int
		setupTokens    int
		olderThan      time.Duration
		expectedDelete int64
	}{
		{
			name:           "success - deletes old revoked tokens",
			loginSuffix:    "delrev1",
			setupTokens:    3,
			revokeIndexes:  []int{0, 1},
			olderThan:      1 * time.Hour,
			expectedDelete: 2,
		},
		{
			name:           "success - does not delete active tokens",
			loginSuffix:    "delrev2",
			setupTokens:    3,
			revokeIndexes:  []int{},
			olderThan:      1 * time.Hour,
			expectedDelete: 0,
		},
		{
			name:           "success - does not delete recent revoked tokens",
			loginSuffix:    "delrev3",
			setupTokens:    2,
			revokeIndexes:  []int{0, 1},
			olderThan:      -1 * time.Hour,
			expectedDelete: 0,
		},
		{
			name:           "success - no tokens to delete",
			loginSuffix:    "delrev4",
			setupTokens:    0,
			revokeIndexes:  []int{},
			olderThan:      -2 * time.Hour,
			expectedDelete: 0,
		},
		{
			name:           "success - mixed revoked and active",
			loginSuffix:    "delrev5",
			setupTokens:    4,
			revokeIndexes:  []int{0, 2},
			olderThan:      -3 * time.Hour,
			expectedDelete: 0,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			// Setup.
			user := s.createTestUser(ctx, "delrevuser_"+tc.loginSuffix)

			tokenIDs := make([]int64, 0, tc.setupTokens)
			for i := range tc.setupTokens {
				tokenHash := s.generateTokenHash(tc.loginSuffix, i)
				insertedToken := s.createTestToken(ctx, user.ID, tokenHash, time.Now().Add(24*time.Hour))
				tokenIDs = append(tokenIDs, insertedToken.ID)
			}

			for _, idx := range tc.revokeIndexes {
				err := s.repo.RevokeRefreshTokenByID(ctx, tokenIDs[idx])
				s.Require().NoError(err)
			}

			// Test.
			deletedCount, err := s.repo.DeleteRevokedTokens(ctx, time.Now().Add(tc.olderThan))

			// Check.
			s.Require().NoError(err)
			s.Equal(tc.expectedDelete, deletedCount)
		})
	}
}
