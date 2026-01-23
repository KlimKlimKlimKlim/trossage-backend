package repository

import (
	"context"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/model"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres/repository/query/message"
)

func (r *Repository) CreateMessage(ctx context.Context, chatID, senderID int64, text string) (model.Message, error) {
	var msg model.Message

	err := r.db.QueryRow(ctx, message.CreateMessageQuery, chatID, senderID, text).Scan(
		&msg.ID,
		&msg.ChatID,
		&msg.SenderID,
		&msg.Text,
		&msg.CreatedAt,
	)
	if err != nil {
		return msg, err
	}

	return msg, nil
}

func (r *Repository) SelectMessages(ctx context.Context, chatID int64, limit, offset int) ([]model.Message, error) {
	rows, err := r.db.Query(ctx, message.SelectMessagesQuery, chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Message, 0, limit)

	for rows.Next() {
		var message model.Message

		err = rows.Scan(
			&message.ID,
			&message.ChatID,
			&message.SenderID,
			&message.Text,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, message)
	}

	return result, rows.Err()
}

func (r *Repository) CountChatMessages(ctx context.Context, chatID int64) (int, error) {
	var count int

	err := r.db.QueryRow(ctx, message.CountChatMessagesQuery, chatID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
