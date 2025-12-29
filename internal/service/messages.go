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
) (models.Message, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return models.Message{}, derrors.ErrMessageIsEmpty
	}

	var message models.Message

	err := s.RepoManager.InTx(ctx, func(tx postgres.IRepository) error {
		_, err := tx.SelectChatByID(ctx, chatID)
		if err != nil {
			return fmt.Errorf("failed to select chat: %w", err)
		}

		_, err = tx.SelectUserByID(ctx, senderID)
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
		return models.Message{}, err
	}

	return message, nil
}

func (s *Service) GetMessages(
	ctx context.Context,
	userID, chatID int64,
	limit, offset int,
) ([]models.Message, int, error) {
	var (
		messages []models.Message
		total    int
	)

	err := s.RepoManager.InTx(ctx, func(tx postgres.IRepository) error {
		isMember, err := tx.IsUserMember(ctx, chatID, userID)
		if err != nil {
			return fmt.Errorf("failed to check is user member: %w", err)
		}

		if !isMember {
			return derrors.ErrUserIsNotMember
		}

		messages, err = tx.SelectMessages(ctx, chatID, limit, offset)
		if err != nil {
			return fmt.Errorf("failed to select messages: %w", err)
		}

		total, err = tx.CountChatMessages(ctx, chatID)
		if err != nil {
			return fmt.Errorf("failed to count messages: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}
