package service

import (
	"context"
	"fmt"
	"strings"

	derrors "github.com/GlaciemArgentum/trossage-backend/internal/errors"
	"github.com/GlaciemArgentum/trossage-backend/internal/model"
	"github.com/GlaciemArgentum/trossage-backend/internal/websocket/hub"
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

	err := s.InTx(ctx, func(txS *Service) error {
		err := txS.isUserMember(ctx, chatID, senderID)
		if err != nil {
			return err
		}

		message, err = txS.Repo.CreateMessage(ctx, chatID, senderID, text)
		if err != nil {
			return fmt.Errorf("failed to create message: %w", err)
		}

		userIDs, err = txS.Repo.SelectChatMembers(ctx, chatID)
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

	err := s.InReadOnlyTx(ctx, func(txS *Service) error {
		err := txS.isUserMember(ctx, chatID, userID)
		if err != nil {
			return err
		}

		messages, err = txS.Repo.SelectMessages(ctx, chatID, limit, offset)
		if err != nil {
			return fmt.Errorf("failed to select messages: %w", err)
		}

		total, err = txS.Repo.CountChatMessages(ctx, chatID)
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

func (s *Service) isUserMember(ctx context.Context, chatID, userID int64) error {
	isMember, err := s.Repo.IsUserMember(ctx, chatID, userID)
	if err != nil {
		return fmt.Errorf("failed to check is user member: %w", err)
	}

	if !isMember {
		return derrors.ErrUserIsNotMember
	}

	return nil
}
