package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/GlaciemArgentum/trossage-backend/internal/http/dto"
	"github.com/GlaciemArgentum/trossage-backend/internal/websocket/hub"
)

func (s *Service) SendTyping(ctx context.Context, senderID, chatID int64, typing dto.TypingUpdateRequest) error {
	var userIDs []int64

	err := s.InReadOnlyTx(ctx, func(txS *Service) error {
		err := txS.isUserMember(ctx, chatID, senderID)
		if err != nil {
			return err
		}

		userIDs, err = txS.Repo.SelectChatMembers(ctx, chatID)
		if err != nil {
			return fmt.Errorf("failed to select chat members: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	userIDs = slices.DeleteFunc(userIDs, func(id int64) bool {
		return id == senderID
	})

	s.WSHub.BroadcastToUsers(userIDs, hub.NewTypingEvent(senderID, chatID, typing.Operations))

	return nil
}
