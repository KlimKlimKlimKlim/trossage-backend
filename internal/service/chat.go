package service

import (
	"context"
	"errors"
	"fmt"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/model"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/websocket/hub"
)

func (s *Service) CreateChat(ctx context.Context, userID, otherUserID int64) (model.Chat, model.User, error) {
	if otherUserID == userID {
		return model.Chat{}, model.User{}, derrors.ErrCannotChatWithYourself
	}

	var (
		chat        model.Chat
		currentUser model.User
		otherUser   model.User
	)

	err := s.InTx(ctx, func(txS *Service) error {
		var err error

		currentUser, err = txS.Repo.SelectUserByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to select user: %w", err)
		}

		otherUser, err = txS.Repo.SelectUserByID(ctx, otherUserID)
		if err != nil {
			return fmt.Errorf("failed to select other user: %w", err)
		}

		existingChatID, err := txS.Repo.SelectChatBetweenUsers(ctx, userID, otherUserID)
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

		chat, err = txS.Repo.InsertChat(ctx, model.ChatTypePrivate)
		if err != nil {
			return fmt.Errorf("failed to insert chat: %w", err)
		}

		err = txS.Repo.InsertChatParticipants(ctx, chat.ID, userID, otherUserID)
		if err != nil {
			return fmt.Errorf("failed to insert participants: %w", err)
		}

		return nil
	})
	if err != nil {
		return model.Chat{}, model.User{}, err
	}

	s.WSHub.BroadcastToUser(otherUserID, hub.NewChatCreatedEvent(chat, currentUser))

	return chat, otherUser, nil
}

func (s *Service) GetUserChats(
	ctx context.Context,
	userID int64,
	limit, offset int,
) ([]model.ChatWithDetails, int, error) {
	var (
		chats []model.ChatWithDetails
		total int
	)

	err := s.InReadOnlyTx(ctx, func(txS *Service) error {
		var err error

		chats, err = txS.Repo.SelectUserChats(ctx, userID, limit, offset)
		if err != nil {
			return fmt.Errorf("failed to get user chats: %w", err)
		}

		total, err = txS.Repo.CountUserChats(ctx, userID)
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
