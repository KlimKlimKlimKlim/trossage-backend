package service

import (
	"context"
	"errors"
	"fmt"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres"
)

func (s *Service) createTokens(ctx context.Context, tx postgres.IRepository, userID int64) (string, string, error) {
	accessString, err := s.AccessJWTController.GenerateSignedToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshString, storeToken, err := s.RefreshJWTController.GenerateSignedTokenAndModel(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	if _, err = tx.InsertRefreshToken(ctx, storeToken); err != nil {
		return "", "", fmt.Errorf("failed to insert token: %w", err)
	}

	return accessString, refreshString, nil
}

func (s *Service) RefreshToken(ctx context.Context, userID, oldTokenID int64) (string, string, error) {
	var accessString, refreshString string

	err := s.RepoManager.InTx(ctx, func(tx postgres.IRepository) error {
		user, err := tx.SelectAuthUserByID(ctx, userID)
		if err != nil {
			if errors.Is(err, derrors.ErrUserNotFound) {
				return fmt.Errorf("user not found: %w", derrors.ErrUnauthorized)
			}

			return fmt.Errorf("failed to select user: %w", err)
		}

		if user.IsDeleted() {
			return fmt.Errorf("user is deleted: %w", derrors.ErrUnauthorized)
		}

		if err = tx.RevokeRefreshTokenByID(ctx, oldTokenID); err != nil {
			if errors.Is(err, derrors.ErrTokenNotFound) {
				return fmt.Errorf("token not found: %w", derrors.ErrUnauthorized)
			}

			return fmt.Errorf("failed to revoke old token: %w", err)
		}

		accessString, refreshString, err = s.createTokens(ctx, tx, userID)
		if err != nil {
			return fmt.Errorf("failed to create new tokens: %w", err)
		}

		return nil
	})
	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}

func (s *Service) Logout(ctx context.Context, tokenID int64) error {
	if err := s.RepoManager.Repo().RevokeRefreshTokenByID(ctx, tokenID); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}

func (s *Service) LogoutAll(ctx context.Context, userID int64) error {
	if err := s.RepoManager.Repo().RevokeRefreshTokensByUserID(ctx, userID); err != nil {
		return fmt.Errorf("failed to revoke tokens: %w", err)
	}

	s.WSHub.DisconnectUser(userID)

	return nil
}
