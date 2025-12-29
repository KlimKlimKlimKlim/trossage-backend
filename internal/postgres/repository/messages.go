package repository

import (
	"context"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres/repository/queries/messages"
)

func (r *Repository) CreateMessage(ctx context.Context, chatID, senderID int64, text string) (models.Message, error) {
	var msg models.Message

	err := r.db.QueryRow(ctx, messages.CreateMessageQuery, chatID, senderID, text).Scan(
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
