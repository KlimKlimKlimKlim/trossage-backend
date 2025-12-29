package service

import (
	"context"
	"fmt"
	"strings"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres"
)

func (s *Service) CreateMessage(
	ctx context.Context,
	senderID, chatID int64,
	text string,
) (models.Message, models.User, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return models.Message{}, models.User{}, derrors.ErrMessageIsEmpty
	}

	var (
		message models.Message
		sender  models.User
	)

	err := s.RepoManager.InTx(ctx, func(tx postgres.IRepository) error {
		_, err := tx.SelectChatByID(ctx, chatID)
		if err != nil {
			return fmt.Errorf("failed to select chat: %w", err)
		}

		sender, err = tx.SelectUserByID(ctx, senderID)
		if err != nil {
			return fmt.Errorf("failed to select user: %w", err)
		}

		isMember, err := tx.IsUserMember(ctx, chatID, senderID)
		if err != nil {
			return fmt.Errorf("failed to check is user member: %w", err)
		}

		if !isMember {
			return derrors.ErrUserIsNotMember
		}

		message, err = tx.CreateMessage(ctx, chatID, senderID, text)
		if err != nil {
			return fmt.Errorf("failed to create message: %w", err)
		}

		return nil
	})
	if err != nil {
		return models.Message{}, models.User{}, err
	}

	return message, sender, nil
}
