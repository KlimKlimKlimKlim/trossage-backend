package controller

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/controller/validate"
	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres"
)

func (c *Controller) CreateUser(ctx context.Context, login, password, displayName string) (models.User, string, string, error) {
	login = strings.ToLower(login)
	if !validate.Login(login) {
		return models.User{}, "", "", derrors.ErrInvalidLogin
	}

	if !validate.Password(password) {
		return models.User{}, "", "", derrors.ErrInvalidPassword
	}

	if !validate.DisplayName(displayName) {
		return models.User{}, "", "", derrors.ErrInvalidDisplayName
	}

	ph, err := c.hasher.HashPassword(password)
	if err != nil {
		return models.User{}, "", "", fmt.Errorf("failed to hash password: %w", err)
	}

	var (
		accessString, refreshString string

		user = models.User{
			Login:        login,
			PasswordHash: ph,
			DisplayName:  displayName,
		}
	)

	err = c.RepoManager.InTx(ctx, func(tx postgres.IRepository) error {
		user, err = tx.InsertUser(ctx, user)
		if err != nil {
			return fmt.Errorf("failed to insert user: %w", err)
		}

		accessString, refreshString, err = c.createTokens(ctx, tx, user.ID)
		if err != nil {
			return fmt.Errorf("failed to create tokens: %w", err)
		}

		return nil
	})
	if err != nil {
		return models.User{}, "", "", err
	}

	return user, accessString, refreshString, nil
}

func (c *Controller) LoginUser(ctx context.Context, login, password string) (models.User, string, string, error) {
	var (
		user                        models.User
		accessString, refreshString string
	)

	err := c.RepoManager.InTx(ctx, func(tx postgres.IRepository) error {
		var err error
		user, err = tx.SelectUserByLogin(ctx, login)
		if err != nil {
			if errors.Is(err, derrors.ErrUserNotFound) {
				return fmt.Errorf("user not found: %w", derrors.ErrUnauthorized)
			}

			return fmt.Errorf("failed to select user: %w", err)
		}

		ok, err := c.hasher.VerifyPassword(password, user.PasswordHash)
		if err != nil {
			return fmt.Errorf("failed to verify password: %w", err)
		}

		if !ok {
			return fmt.Errorf("password is wrong: %w", derrors.ErrUnauthorized)
		}

		accessString, refreshString, err = c.createTokens(ctx, tx, user.ID)
		if err != nil {
			return fmt.Errorf("failed to create tokens: %w", err)
		}

		return nil
	})
	if err != nil {
		return models.User{}, "", "", err
	}

	return user, accessString, refreshString, nil
}

func (c *Controller) GetUserByID(ctx context.Context, userID int64) (models.User, error) {
	user, err := c.RepoManager.Repo().SelectUserByID(ctx, userID)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (c *Controller) UpdateUser(ctx context.Context, userID int64, displayName, oldPassword, newPassword string) (models.User, bool, string, error) {
	if !validate.Password(newPassword) {
		return models.User{}, false, "", derrors.ErrInvalidPassword
	}

	if !validate.DisplayName(displayName) {
		return models.User{}, false, "", derrors.ErrInvalidDisplayName
	}

	var (
		updatedUser         models.User
		tokensRevoked       bool
		tokensRevokedReason string
	)

	err := c.RepoManager.InTx(ctx, func(tx postgres.IRepository) error {
		user, err := tx.SelectUserByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		if user.IsDeleted() {
			return fmt.Errorf("user is deleted: %w", derrors.ErrUserNotFound)
		}

		if newPassword != "" {
			ok, err := c.hasher.VerifyPassword(oldPassword, user.PasswordHash)
			if err != nil {
				return fmt.Errorf("failed to verify password: %w", err)
			}

			if !ok {
				return derrors.ErrUnauthorized
			}

			isSame, err := c.hasher.VerifyPassword(newPassword, user.PasswordHash)
			if err != nil {
				return fmt.Errorf("failed to verify password: %w", err)
			}

			if isSame {
				return derrors.ErrSamePassword
			}

			newHash, err := c.hasher.HashPassword(newPassword)
			if err != nil {
				return fmt.Errorf("failed to hash password: %w", err)
			}

			user.PasswordHash = newHash
			tokensRevoked = true
			tokensRevokedReason = models.PasswordChangedReason

			if err = tx.RevokeRefreshTokensByUserID(ctx, userID); err != nil {
				return fmt.Errorf("failed to revoke tokens: %w", err)
			}
		}

		if displayName != "" {
			user.DisplayName = displayName
		}

		updatedUser, err = tx.UpdateUser(ctx, user)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		return nil
	})
	if err != nil {
		return models.User{}, false, "", err
	}

	return updatedUser, tokensRevoked, tokensRevokedReason, nil
}

func (c *Controller) DeleteUser(ctx context.Context, userID int64, password string) error {
	err := c.RepoManager.InTx(ctx, func(tx postgres.IRepository) error {
		user, err := tx.SelectUserByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		ok, err := c.hasher.VerifyPassword(password, user.PasswordHash)
		if err != nil {
			return fmt.Errorf("failed to verify password: %w", err)
		}

		if !ok {
			return derrors.ErrUnauthorized
		}

		if user.IsDeleted() {
			return fmt.Errorf("user is deleted: %w", derrors.ErrUserNotFound)
		}

		if err = tx.DeleteUser(ctx, userID); err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}

		return nil
	})

	return err
}
