package dto

import (
	"time"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/model"
)

type SendMessageRequest struct {
	Text string `json:"text" binding:"required,min=1,max=4096" example:"Привет!"`
}

type MessageResponse struct {
	CreatedAt time.Time `json:"created_at" example:"2025-12-29T13:00:00Z"`
	Text      string    `json:"text"       example:"Привет!"`
	SenderID  int64     `json:"sender_id"  example:"789"`
	ID        int64     `json:"id"         example:"456"`
	ChatID    int64     `json:"chat_id"    example:"123"`
}

func (dto *MessageResponse) Fill(msg model.Message) {
	dto.ID = msg.ID
	dto.ChatID = msg.ChatID
	dto.Text = msg.Text
	dto.CreatedAt = msg.CreatedAt
	dto.SenderID = msg.SenderID
}

type MessagesResponse struct {
	Messages []MessageResponse `json:"messages"`
	Total    int               `json:"total"    example:"150"`
	Limit    int               `json:"limit"    example:"20"`
	Offset   int               `json:"offset"   example:"0"`
}

func (dto *MessagesResponse) Fill(messages []model.Message, limit, offset, total int) {
	dto.Limit = limit
	dto.Offset = offset
	dto.Total = total

	dto.Messages = make([]MessageResponse, len(messages))
	for i, message := range messages {
		dto.Messages[i].Fill(message)
	}
}
