package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	derrors "github.com/GlaciemArgentum/trossage-backend/internal/errors"
	"github.com/GlaciemArgentum/trossage-backend/internal/model"
	"github.com/GlaciemArgentum/trossage-backend/internal/service/validate"
)

const (
	PasswordChangedReason = "password changed"
)

func (s *Service) CreateUser(
	ctx context.Context,
	login, password, displayName string,
) (model.User, model.JWTPair, error) {
	login = strings.ToLower(login)
	if !validate.Login(login) {
		return model.User{}, model.JWTPair{}, derrors.ErrInvalidLogin
	}

	if !validate.Password(password) {
		return model.User{}, model.JWTPair{}, derrors.ErrInvalidPassword
	}

	if !validate.DisplayName(displayName) {
		return model.User{}, model.JWTPair{}, derrors.ErrInvalidDisplayName
	}

	ph, err := s.hasher.HashPassword(password)
	if err != nil {
		return model.User{}, model.JWTPair{}, fmt.Errorf("failed to hash password: %w", err)
	}

	var (
		jwtPair model.JWTPair
		user    model.AuthUser
	)

	user.Login = login
	user.DisplayName = displayName
	user.PasswordHash = ph

	err = s.InTx(ctx, func(txS *Service) error {
		user, err = txS.Repo.InsertUser(ctx, user)
		if err != nil {
			return fmt.Errorf("failed to insert user: %w", err)
		}

		jwtPair.AccessToken, jwtPair.RefreshToken, err = txS.createTokens(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("failed to create tokens: %w", err)
		}

		return nil
	})
	if err != nil {
		return model.User{}, model.JWTPair{}, err
	}

	return user.User, jwtPair, nil
}

func (s *Service) LoginUser(ctx context.Context, login, password string) (model.User, model.JWTPair, error) {
	user, err := s.Repo.SelectAuthUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, derrors.ErrUserNotFound) {
			return model.User{}, model.JWTPair{}, fmt.Errorf("user not found: %w", derrors.ErrUnauthorized)
		}

		return model.User{}, model.JWTPair{}, fmt.Errorf("failed to select user: %w", err)
	}

	ok, err := s.hasher.VerifyPassword(password, user.PasswordHash)
	if err != nil {
		return model.User{}, model.JWTPair{}, fmt.Errorf("failed to verify password: %w", err)
	}

	if !ok {
		return model.User{}, model.JWTPair{}, fmt.Errorf("password is wrong: %w", derrors.ErrUnauthorized)
	}

	var JWTPair model.JWTPair

	JWTPair.AccessToken, JWTPair.RefreshToken, err = s.createTokens(ctx, user.ID)
	if err != nil {
		return model.User{}, model.JWTPair{}, fmt.Errorf("failed to create tokens: %w", err)
	}

	return user.User, JWTPair, nil
}

func (s *Service) GetUserByID(ctx context.Context, userID int64) (model.User, error) {
	user, err := s.Repo.SelectUserByID(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *Service) UpdateUser(
	ctx context.Context,
	userID int64,
	displayName, oldPassword, newPassword string,
) (model.User, model.TokenRevocation, error) {
	var (
		updatedUser     model.AuthUser
		tokenRevocation model.TokenRevocation
	)

	err := s.InTx(ctx, func(txS *Service) error {
		user, err := txS.Repo.SelectAuthUserByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		if newPassword != "" {
			if err = txS.updateUserPassword(ctx, &user, oldPassword, newPassword); err != nil {
				return err
			}

			tokenRevocation.Revoked = true
			tokenRevocation.Reason = PasswordChangedReason
		}

		if displayName != "" {
			if err = txS.updateUserDisplayName(&user, displayName); err != nil {
				return err
			}
		}

		updatedUser, err = txS.Repo.UpdateUser(ctx, user)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		return nil
	})
	if err != nil {
		return model.User{}, model.TokenRevocation{}, err
	}

	if newPassword != "" {
		s.WSHub.DisconnectUser(userID)
	}

	return updatedUser.User, tokenRevocation, nil
}

func (s *Service) updateUserPassword(
	ctx context.Context,
	user *model.AuthUser,
	oldPassword, newPassword string,
) error {
	if !validate.Password(newPassword) {
		return derrors.ErrInvalidPassword
	}

	ok, err := s.hasher.VerifyPassword(oldPassword, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to verify password: %w", err)
	}

	if !ok {
		return derrors.ErrUnauthorized
	}

	isSame, err := s.hasher.VerifyPassword(newPassword, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to verify password: %w", err)
	}

	if isSame {
		return derrors.ErrSamePassword
	}

	newHash, err := s.hasher.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = newHash

	if err = s.Repo.RevokeRefreshTokensByUserID(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to revoke tokens: %w", err)
	}

	return nil
}

func (s *Service) updateUserDisplayName(user *model.AuthUser, displayName string) error {
	if !validate.DisplayName(displayName) {
		return derrors.ErrInvalidDisplayName
	}

	user.DisplayName = displayName

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, userID int64, password string) error {
	err := s.InTx(ctx, func(txS *Service) error {
		user, err := txS.Repo.SelectAuthUserByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		ok, err := txS.hasher.VerifyPassword(password, user.PasswordHash)
		if err != nil {
			return fmt.Errorf("failed to verify password: %w", err)
		}

		if !ok {
			return derrors.ErrUnauthorized
		}

		if err = txS.Repo.RevokeRefreshTokensByUserID(ctx, userID); err != nil {
			return fmt.Errorf("failed to revoke tokens: %w", err)
		}

		if err = txS.Repo.DeleteUser(ctx, userID); err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	s.WSHub.DisconnectUser(userID)

	return nil
}

func (s *Service) SearchUsersByLogin(
	ctx context.Context,
	userID int64,
	prefix string,
	limit, offset int,
) ([]model.User, int, error) {
	var (
		users []model.User
		total int
	)

	err := s.InReadOnlyTx(ctx, func(txS *Service) error {
		var err error

		users, err = txS.Repo.SelectUsersByLoginPrefix(ctx, userID, prefix, limit, offset)
		if err != nil {
			return fmt.Errorf("failed to search users: %w", err)
		}

		total, err = txS.Repo.CountUsersByLoginPrefix(ctx, userID, prefix)
		if err != nil {
			return fmt.Errorf("failed to count users: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
