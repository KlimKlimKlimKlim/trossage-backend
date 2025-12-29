package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres/repository/queries/chats"
)

func (r *Repository) SelectChatBetweenUsers(ctx context.Context, userID1, userID2 int64) (int64, error) {
	var chatID int64

	err := r.db.QueryRow(ctx, chats.SelectChatBetweenUsersQuery, userID1, userID2).Scan(&chatID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, derrors.ErrChatNotFound
		}

		return 0, err
	}

	return chatID, nil
}

func (r *Repository) InsertChat(ctx context.Context, chatType models.ChatType) (models.Chat, error) {
	var chat models.Chat

	err := r.db.QueryRow(ctx, chats.InsertChatQuery, chatType).Scan(
		&chat.ID,
		&chat.Type,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)
	if err != nil {
		return models.Chat{}, err
	}

	return chat, nil
}

func (r *Repository) InsertChatParticipants(ctx context.Context, chatID int64, userIDs ...int64) error {
	if len(userIDs) == 0 {
		return derrors.ErrEmptyInput
	}

	_, err := r.db.Exec(ctx, chats.InsertChatParticipantsQuery, chatID, userIDs)

	return err
}
