package dto

import (
	"time"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
)

type SendMessageRequest struct {
	Text string `json:"text" binding:"required,min=1,max=4096" example:"Привет!"`
}

type MessageResponse struct {
	CreatedAt time.Time    `json:"created_at" example:"2025-12-29T13:00:00Z"`
	Text      string       `json:"text"       example:"Привет!"`
	Sender    UserResponse `json:"sender"`
	ID        int64        `json:"id"         example:"456"`
	ChatID    int64        `json:"chat_id"    example:"123"`
}

func (dto *MessageResponse) Fill(msg models.Message, sender models.User) {
	dto.ID = msg.ID
	dto.ChatID = msg.ChatID
	dto.Text = msg.Text
	dto.CreatedAt = msg.CreatedAt
	dto.Sender.Fill(sender)
}
