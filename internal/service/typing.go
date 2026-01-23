package service

import (
	"context"
	"fmt"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/dto"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/websocket/hub"
)

func (s *Service) SendTyping(ctx context.Context, senderID, chatID int64, typing dto.TypingUpdateRequest) error {
	var userIDs []int64

	err := s.RepoManager.InTx(ctx, func(tx postgres.IRepository) error {
		isMember, err := tx.IsUserMember(ctx, chatID, senderID)
		if err != nil {
			return fmt.Errorf("failed to check is user member: %w", err)
		}

		if !isMember {
			return derrors.ErrUserIsNotMember
		}

		userIDs, err = tx.SelectChatMembers(ctx, chatID)
		if err != nil {
			return fmt.Errorf("failed to select chat members: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	s.WSHub.BroadcastToUsers(userIDs, hub.NewTypingEvent(senderID, chatID, typing.Operations))

	return nil
}
