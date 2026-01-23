package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/model"
	query "github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres/repository/query/chat"
)

func (r *Repository) SelectChatBetweenUsers(ctx context.Context, userID1, userID2 int64) (int64, error) {
	var chatID int64

	err := r.db.QueryRow(ctx, query.SelectChatBetweenUsersQuery, userID1, userID2).Scan(&chatID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, derrors.ErrChatNotFound
		}

		return 0, err
	}

	return chatID, nil
}

func (r *Repository) InsertChat(ctx context.Context, chatType model.ChatType) (model.Chat, error) {
	var chat model.Chat

	err := r.db.QueryRow(ctx, query.InsertChatQuery, chatType).Scan(
		&chat.ID,
		&chat.Type,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)
	if err != nil {
		return model.Chat{}, err
	}

	return chat, nil
}

func (r *Repository) InsertChatParticipants(ctx context.Context, chatID int64, userIDs ...int64) error {
	if len(userIDs) == 0 {
		return derrors.ErrEmptyInput
	}

	_, err := r.db.Exec(ctx, query.InsertChatParticipantsQuery, chatID, userIDs)

	return err
}

func (r *Repository) SelectUserChats(
	ctx context.Context,
	userID int64,
	limit, offset int,
) ([]model.ChatWithDetails, error) {
	rows, err := r.db.Query(ctx, query.SelectUserChatsQuery, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.ChatWithDetails, 0, limit)

	for rows.Next() {
		chat, err := r.scanChatWithDetails(rows)
		if err != nil {
			return nil, err
		}

		result = append(result, chat)
	}

	return result, rows.Err()
}

func (r *Repository) CountUserChats(ctx context.Context, userID int64) (int, error) {
	var count int

	err := r.db.QueryRow(ctx, query.CountUserChatsQuery, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) IsUserMember(ctx context.Context, chatID, userID int64) (bool, error) {
	var isMember bool

	err := r.db.QueryRow(ctx, query.IsUserMemberQuery, chatID, userID).Scan(&isMember)
	if err != nil {
		return false, err
	}

	return isMember, nil
}

func (r *Repository) SelectChatByID(ctx context.Context, chatID int64) (model.Chat, error) {
	var chat model.Chat

	err := r.db.QueryRow(ctx, query.SelectChatByIDQuery, chatID).Scan(
		&chat.ID,
		&chat.Type,
		&chat.CreatedAt,
		&chat.UpdatedAt,
		&chat.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Chat{}, derrors.ErrChatNotFound
		}

		return model.Chat{}, err
	}

	return chat, nil
}

func (r *Repository) SelectChatMembers(ctx context.Context, chatID int64) ([]int64, error) {
	rows, err := r.db.Query(ctx, query.SelectChatMembersQuery, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userIDs := make([]int64, 0)

	for rows.Next() {
		var userID int64
		if err = rows.Scan(&userID); err != nil {
			return nil, err
		}

		userIDs = append(userIDs, userID)
	}

	return userIDs, rows.Err()
}
