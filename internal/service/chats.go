package service

import (
	"context"
	"errors"
	"fmt"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres"
)

func (s *Service) CreateChat(ctx context.Context, userID, otherUserID int64) (models.Chat, models.User, error) {
	if otherUserID == userID {
		return models.Chat{}, models.User{}, derrors.ErrCannotChatWithYourself
	}

	var (
		chat      models.Chat
		otherUser models.User
	)

	err := s.RepoManager.InTx(ctx, func(tx postgres.IRepository) error {
		_, err := s.getAndCheckUserByID(ctx, tx, userID)
		if err != nil {
			return fmt.Errorf("failed to select user: %w", err)
		}

		otherUser, err = s.getAndCheckUserByID(ctx, tx, otherUserID)
		if err != nil {
			return fmt.Errorf("failed to select other user: %w", err)
		}

		existingChatID, err := tx.SelectChatBetweenUsers(ctx, userID, otherUserID)
		if err != nil && !errors.Is(err, derrors.ErrChatNotFound) {
			return fmt.Errorf("failed to check existing chat: %w", err)
		}

		if existingChatID != 0 {
			return fmt.Errorf(
				"chat already exists with id = %d: %w",
				existingChatID,
				derrors.ErrChatAlreadyExists,
			)
		}

		chat, err = tx.InsertChat(ctx, models.ChatTypePrivate)
		if err != nil {
			return fmt.Errorf("failed to insert chat: %w", err)
		}

		err = tx.InsertChatParticipants(ctx, chat.ID, userID, otherUserID)
		if err != nil {
			return fmt.Errorf("failed to insert participants: %w", err)
		}

		return nil
	})
	if err != nil {
		return models.Chat{}, models.User{}, err
	}

	return chat, otherUser, nil
}

func (s *Service) GetUserChats(
	ctx context.Context,
	userID int64,
	limit, offset int,
) ([]models.ChatWithDetails, int, error) {
	var (
		chats []models.ChatWithDetails
		total int
	)

	err := s.RepoManager.InTx(ctx, func(tx postgres.IRepository) error {
		var err error

		chats, err = tx.SelectUserChats(ctx, userID, limit, offset)
		if err != nil {
			return fmt.Errorf("failed to get user chats: %w", err)
		}

		total, err = tx.CountUserChats(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to count user chats: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return chats, total, nil
}
