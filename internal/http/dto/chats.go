package dto

import (
	"time"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
)

type CreateChatRequest struct {
	UserID int64 `json:"user_id" binding:"required,min=1" example:"456"`
}

type ChatResponse struct {
	CreatedAt time.Time     `json:"created_at"           example:"2025-12-23T00:00:00Z"`
	OtherUser *UserResponse `json:"other_user,omitempty"`
	Type      string        `json:"type"                 example:"private"`
	ID        int64         `json:"id"                   example:"123"`
}

func (dto *ChatResponse) Fill(chat models.Chat, otherUser models.User) {
	dto.ID = chat.ID
	dto.Type = string(chat.Type)
	dto.CreatedAt = chat.CreatedAt

	if chat.Type == models.ChatTypePrivate && otherUser.ID != 0 {
		other := &UserResponse{}
		other.Fill(otherUser)
		dto.OtherUser = other
	}
}
