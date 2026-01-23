package service

import (
	"context"
	"fmt"
	"strings"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/model"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/websocket/hub"
)

func (s *Service) CreateMessage(
	ctx context.Context,
	senderID, chatID int64,
	text string,
) (model.Message, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return model.Message{}, derrors.ErrMessageIsEmpty
	}

	var (
		message model.Message
		userIDs []int64
	)

	err := s.RepoManager.InTx(ctx, func(tx postgres.IRepository) error {
		err := s.isUserMember(ctx, tx, chatID, senderID)
		if err != nil {
			return err
		}

		message, err = tx.CreateMessage(ctx, chatID, senderID, text)
		if err != nil {
			return fmt.Errorf("failed to create message: %w", err)
		}

		userIDs, err = tx.SelectChatMembers(ctx, chatID)
		if err != nil {
			return fmt.Errorf("failed to select chat members: %w", err)
		}

		return nil
	})
	if err != nil {
		return model.Message{}, err
	}

	s.WSHub.BroadcastToUsers(userIDs, hub.NewNewMessageEvent(message))

	return message, nil
}

func (s *Service) GetMessages(
	ctx context.Context,
	userID, chatID int64,
	limit, offset int,
) ([]model.Message, int, error) {
	var (
		messages []model.Message
		total    int
	)

	err := s.RepoManager.InTx(ctx, func(tx postgres.IRepository) error {
		err := s.isUserMember(ctx, tx, chatID, userID)
		if err != nil {
			return err
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

func (s *Service) isUserMember(ctx context.Context, tx postgres.IRepository, chatID, userID int64) error {
	isMember, err := tx.IsUserMember(ctx, chatID, userID)
	if err != nil {
		return fmt.Errorf("failed to check is user member: %w", err)
	}

	if !isMember {
		return derrors.ErrUserIsNotMember
	}

	return nil
}
