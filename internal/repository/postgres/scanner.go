package postgres

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"github.com/GlaciemArgentum/trossage-backend/internal/model"
)

type nullableUserFields struct {
	CreatedAt   sql.NullTime
	Login       sql.NullString
	DisplayName sql.NullString
	ID          sql.NullInt64
}

type nullableLastMessageFields struct {
	CreatedAt       sql.NullTime
	SenderCreatedAt sql.NullTime
	Text            sql.NullString
	SenderLogin     sql.NullString
	SenderDisplay   sql.NullString
	SenderID        sql.NullInt64
}

func scanNullableUser(fields nullableUserFields) model.User {
	if !fields.ID.Valid {
		return model.User{}
	}

	return model.User{
		ID:          fields.ID.Int64,
		Login:       fields.Login.String,
		DisplayName: fields.DisplayName.String,
		CreatedAt:   fields.CreatedAt.Time,
	}
}

func scanNullableLastMessage(fields nullableLastMessageFields) model.LastMessage {
	if !fields.SenderID.Valid {
		return model.LastMessage{}
	}

	return model.LastMessage{
		Text:      fields.Text.String,
		CreatedAt: fields.CreatedAt.Time,
		Sender: model.User{
			ID:          fields.SenderID.Int64,
			Login:       fields.SenderLogin.String,
			DisplayName: fields.SenderDisplay.String,
			CreatedAt:   fields.SenderCreatedAt.Time,
		},
	}
}

func (r *Repository) scanChatWithDetails(rows pgx.Rows) (model.ChatWithDetails, error) {
	var (
		chat      model.ChatWithDetails
		otherUser nullableUserFields
		lastMsg   nullableLastMessageFields
	)

	err := rows.Scan(
		&chat.ID,
		&chat.Type,
		&chat.CreatedAt,
		&lastMsg.Text,
		&lastMsg.SenderID,
		&lastMsg.CreatedAt,
		&lastMsg.SenderLogin,
		&lastMsg.SenderDisplay,
		&lastMsg.SenderCreatedAt,
		&otherUser.ID,
		&otherUser.Login,
		&otherUser.DisplayName,
		&otherUser.CreatedAt,
	)
	if err != nil {
		return chat, err
	}

	chat.OtherUser = scanNullableUser(otherUser)
	chat.LastMessage = scanNullableLastMessage(lastMsg)

	return chat, nil
}
