package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/service/validate"
)

func (c *Service) CreateUser(
	ctx context.Context,
	login, password, displayName string,
) (models.User, models.JWTPair, error) {
	login = strings.ToLower(login)
	if !validate.Login(login) {
		return models.User{}, models.JWTPair{}, derrors.ErrInvalidLogin
	}

	if !validate.Password(password) {
		return models.User{}, models.JWTPair{}, derrors.ErrInvalidPassword
	}

	if !validate.DisplayName(displayName) {
		return models.User{}, models.JWTPair{}, derrors.ErrInvalidDisplayName
	}

	ph, err := c.hasher.HashPassword(password)
	if err != nil {
		return models.User{}, models.JWTPair{}, fmt.Errorf("failed to hash password: %w", err)
	}

	var (
		jwtPair models.JWTPair

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

		jwtPair.AccessToken, jwtPair.RefreshToken, err = c.createTokens(ctx, tx, user.ID)
		if err != nil {
			return fmt.Errorf("failed to create tokens: %w", err)
		}

		return nil
	})
	if err != nil {
		return models.User{}, models.JWTPair{}, err
	}

	return user, jwtPair, nil
}

func (c *Service) LoginUser(ctx context.Context, login, password string) (models.User, models.JWTPair, error) {
	var (
		user    models.User
		JWTPair models.JWTPair
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

		JWTPair.AccessToken, JWTPair.RefreshToken, err = c.createTokens(ctx, tx, user.ID)
		if err != nil {
			return fmt.Errorf("failed to create tokens: %w", err)
		}

		return nil
	})
	if err != nil {
		return models.User{}, models.JWTPair{}, err
	}

	return user, JWTPair, nil
}

func (c *Service) GetUserByID(ctx context.Context, userID int64) (models.User, error) {
	user, err := c.RepoManager.Repo().SelectUserByID(ctx, userID)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (c *Service) UpdateUser(
	ctx context.Context,
	userID int64,
	displayName, oldPassword, newPassword string,
) (models.User, models.TokenRevocation, error) {
	var (
		updatedUser     models.User
		tokenRevocation models.TokenRevocation
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
			if err = c.updateUserPassword(ctx, tx, &user, oldPassword, newPassword); err != nil {
				return err
			}

			tokenRevocation.Revoked = true
			tokenRevocation.Reason = models.PasswordChangedReason
		}

		if displayName != "" {
			if err = c.updateUserDisplayName(&user, displayName); err != nil {
				return err
			}
		}

		updatedUser, err = tx.UpdateUser(ctx, user)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		return nil
	})
	if err != nil {
		return models.User{}, models.TokenRevocation{}, err
	}

	return updatedUser, tokenRevocation, nil
}

func (c *Service) updateUserPassword(
	ctx context.Context,
	tx postgres.IRepository,
	user *models.User,
	oldPassword, newPassword string,
) error {
	if !validate.Password(newPassword) {
		return derrors.ErrInvalidPassword
	}

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

	if err = tx.RevokeRefreshTokensByUserID(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to revoke tokens: %w", err)
	}

	return nil
}

func (c *Service) updateUserDisplayName(user *models.User, displayName string) error {
	if !validate.DisplayName(displayName) {
		return derrors.ErrInvalidDisplayName
	}

	user.DisplayName = displayName

	return nil
}

func (c *Service) DeleteUser(ctx context.Context, userID int64, password string) error {
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

func (c *Service) SearchUsersByLogin(
	ctx context.Context,
	query string,
	limit, offset int,
) ([]models.User, int, error) {
	users, err := c.RepoManager.Repo().SelectUsersByLoginPrefix(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}

	total, err := c.RepoManager.Repo().CountUsersByLoginPrefix(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	return users, total, nil
}
