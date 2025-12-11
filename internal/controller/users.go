package controller

import (
	"context"
	"errors"
	"fmt"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres"
)

func (c *Controller) CreateUser(ctx context.Context, login, password, displayName string) (models.User, string, string, error) {
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

	err = c.repoManager.InTx(ctx, func(tx postgres.IRepository) error {
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

func (c *Controller) LoginUser(ctx context.Context, login, password string) (string, string, error) {
	var accessString, refreshString string

	err := c.repoManager.InTx(ctx, func(tx postgres.IRepository) error {
		user, err := tx.SelectUserByLogin(ctx, login)
		if err != nil {
			if errors.Is(err, derrors.ErrUserNotFound) {
				return fmt.Errorf("user not found: %w", derrors.ErrWrongLoginOrPassword)
			}

			return fmt.Errorf("failed to select user: %w", err)
		}

		ok, err := c.hasher.VerifyPassword(password, user.PasswordHash)
		if err != nil {
			return fmt.Errorf("failed to verify password: %w", err)
		}

		if !ok {
			return fmt.Errorf("password is wrong: %w", derrors.ErrWrongLoginOrPassword)
		}

		accessString, refreshString, err = c.createTokens(ctx, tx, user.ID)
		if err != nil {
			return fmt.Errorf("failed to create tokens: %w", err)
		}

		return nil
	})
	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}
