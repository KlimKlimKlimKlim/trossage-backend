package tokencleanup

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/GlaciemArgentum/trossage-backend/internal/config"
	"github.com/GlaciemArgentum/trossage-backend/internal/mocks"
)

var (
	errDatabaseFailure = errors.New("database connection failed")
	errQueryTimeout    = errors.New("query timeout exceeded")
)

func TestCleanup(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		expiredDeleteError  error
		revokedDeleteError  error
		config              *config.TokenCleanupConfig
		name                string
		expectedErrContains string
		expiredDeletedCount int64
		revokedDeletedCount int64
		expectedErr         bool
	}{
		{
			name: "success - deletes expired and revoked tokens",
			config: &config.TokenCleanupConfig{
				Interval:         15 * time.Minute,
				ExpiredRetention: 1 * time.Hour,
				RevokedRetention: 2 * time.Hour,
			},
			expiredDeletedCount: 5,
			revokedDeletedCount: 3,
			expectedErr:         false,
		},
		{
			name: "success - no tokens deleted",
			config: &config.TokenCleanupConfig{
				Interval:         15 * time.Minute,
				ExpiredRetention: 1 * time.Hour,
				RevokedRetention: 1 * time.Hour,
			},
			expiredDeletedCount: 0,
			revokedDeletedCount: 0,
			expectedErr:         false,
		},
		{
			name: "error - delete expired tokens fails",
			config: &config.TokenCleanupConfig{
				Interval:         15 * time.Minute,
				ExpiredRetention: 1 * time.Hour,
				RevokedRetention: 1 * time.Hour,
			},
			expiredDeleteError:  errDatabaseFailure,
			expectedErr:         true,
			expectedErrContains: "failed to delete expired tokens",
		},
		{
			name: "error - delete revoked tokens fails",
			config: &config.TokenCleanupConfig{
				Interval:         15 * time.Minute,
				ExpiredRetention: 1 * time.Hour,
				RevokedRetention: 1 * time.Hour,
			},
			expiredDeletedCount: 5,
			revokedDeleteError:  errQueryTimeout,
			expectedErr:         true,
			expectedErrContains: "failed to delete revoked tokens",
		},
		{
			name: "success - with zero retention",
			config: &config.TokenCleanupConfig{
				Interval:         15 * time.Minute,
				ExpiredRetention: 0,
				RevokedRetention: 0,
			},
			expiredDeletedCount: 10,
			revokedDeletedCount: 7,
			expectedErr:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			log := zap.NewNop()
			repo := mocks.NewRepository(t)
			worker := New(log, repo, tc.config)

			ctx := context.Background()

			if tc.expiredDeleteError != nil {
				repo.EXPECT().
					DeleteExpiredTokens(ctx, mock.AnythingOfType("time.Time")).
					Return(int64(0), tc.expiredDeleteError).
					Once()
			} else {
				repo.EXPECT().
					DeleteExpiredTokens(ctx, mock.AnythingOfType("time.Time")).
					Return(tc.expiredDeletedCount, nil).
					Once()
			}

			if tc.expectedErr && tc.expiredDeleteError == nil {
				repo.EXPECT().
					DeleteRevokedTokens(ctx, mock.AnythingOfType("time.Time")).
					Return(int64(0), tc.revokedDeleteError).
					Once()
			} else if !tc.expectedErr {
				repo.EXPECT().
					DeleteRevokedTokens(ctx, mock.AnythingOfType("time.Time")).
					Return(tc.revokedDeletedCount, nil).
					Once()
			}

			// Test.
			err := worker.cleanup(ctx)

			// Check.
			if tc.expectedErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("success - stops on context cancel", func(t *testing.T) {
		t.Parallel()

		log := zap.NewNop()
		repo := mocks.NewRepository(t)

		cfg := &config.TokenCleanupConfig{
			Interval:         15 * time.Minute,
			ExpiredRetention: 1 * time.Hour,
			RevokedRetention: 1 * time.Hour,
		}

		worker := New(log, repo, cfg)

		ctx, cancel := context.WithCancel(context.Background())

		repo.EXPECT().
			DeleteExpiredTokens(mock.Anything, mock.AnythingOfType("time.Time")).
			Return(int64(5), nil).
			Maybe()

		repo.EXPECT().
			DeleteRevokedTokens(mock.Anything, mock.AnythingOfType("time.Time")).
			Return(int64(3), nil).
			Maybe()

		done := make(chan error, 1)

		go func() {
			done <- worker.Run(ctx)
		}()

		cancel()

		select {
		case err := <-done:
			require.NoError(t, err)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("worker did not stop within timeout")
		}
	})

	t.Run("success - continues running after cleanup error", func(t *testing.T) {
		t.Parallel()

		log := zap.NewNop()
		repo := mocks.NewRepository(t)

		cfg := &config.TokenCleanupConfig{
			Interval:         15 * time.Minute,
			ExpiredRetention: 1 * time.Hour,
			RevokedRetention: 1 * time.Hour,
		}

		worker := New(log, repo, cfg)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		callCount := 0

		repo.EXPECT().
			DeleteExpiredTokens(mock.Anything, mock.AnythingOfType("time.Time")).
			RunAndReturn(func(_ context.Context, _ time.Time) (int64, error) {
				callCount++
				return 0, errDatabaseFailure
			}).
			Maybe()

		err := worker.Run(ctx)

		require.NoError(t, err)
		require.GreaterOrEqual(t, callCount, 1, "cleanup should be called at least once")
	})

	t.Run("success - performs initial cleanup immediately", func(t *testing.T) {
		t.Parallel()

		log := zap.NewNop()
		repo := mocks.NewRepository(t)

		cfg := &config.TokenCleanupConfig{
			Interval:         1 * time.Hour,
			ExpiredRetention: 1 * time.Hour,
			RevokedRetention: 1 * time.Hour,
		}

		worker := New(log, repo, cfg)

		ctx, cancel := context.WithCancel(context.Background())

		cleanupCalled := make(chan struct{}, 1)

		repo.EXPECT().
			DeleteExpiredTokens(mock.Anything, mock.AnythingOfType("time.Time")).
			RunAndReturn(func(_ context.Context, _ time.Time) (int64, error) {
				select {
				case cleanupCalled <- struct{}{}:
				default:
				}

				return 5, nil
			}).
			Once()

		repo.EXPECT().
			DeleteRevokedTokens(mock.Anything, mock.AnythingOfType("time.Time")).
			Return(int64(3), nil).
			Once()

		go func() {
			_ = worker.Run(ctx)
		}()

		select {
		case <-cleanupCalled:
			cancel()
		case <-time.After(100 * time.Millisecond):
			cancel()
			t.Fatal("initial cleanup was not called")
		}
	})
}
