package controller

import (
	"context"
	"errors"
	"fmt"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres"
)

func (c *Controller) createTokens(ctx context.Context, tx postgres.IRepository, userID int64) (string, string, error) {
	accessString, err := c.AccessJWTController.GenerateSignedToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshString, storeToken, err := c.RefreshJWTController.GenerateSignedTokenAndModel(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	if _, err = tx.InsertRefreshToken(ctx, storeToken); err != nil {
		return "", "", fmt.Errorf("failed to insert token: %w", err)
	}

	return accessString, refreshString, nil
}

func (c *Controller) RefreshToken(ctx context.Context, userID, oldTokenID int64) (string, string, error) {
	var accessString, refreshString string

	err := c.RepoManager.InTx(ctx, func(tx postgres.IRepository) error {
		if user, err := tx.SelectUserByID(ctx, userID); err != nil {
			switch {
			case errors.Is(err, derrors.ErrUserNotFound):
				return fmt.Errorf("user not found: %w", derrors.ErrUnauthorized)
			case user.IsDeleted():
				return fmt.Errorf("user is deleted: %w", derrors.ErrUnauthorized)
			}

			return fmt.Errorf("failed to select user: %w", err)
		}

		if err := tx.RevokeRefreshTokenByID(ctx, oldTokenID); err != nil {
			if errors.Is(err, derrors.ErrTokenNotFound) {
				return fmt.Errorf("token not found: %w", derrors.ErrUnauthorized)
			}

			return fmt.Errorf("failed to revoke old token: %w", err)
		}

		var err error
		accessString, refreshString, err = c.createTokens(ctx, tx, userID)
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

func (c *Controller) Logout(ctx context.Context, tokenID int64) error {
	if err := c.RepoManager.Repo().RevokeRefreshTokenByID(ctx, tokenID); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}

func (c *Controller) LogoutAll(ctx context.Context, userID int64) error {
	if err := c.RepoManager.Repo().RevokeRefreshTokensByUserID(ctx, userID); err != nil {
		return fmt.Errorf("failed to revoke tokens: %w", err)
	}

	return nil
}
