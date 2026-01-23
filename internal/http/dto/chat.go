package dto

import (
	"time"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/model"
)

type CreateChatRequest struct {
	UserID int64 `json:"user_id" binding:"required,min=1" example:"456"`
}

type ChatResponse struct {
	CreatedAt   time.Time        `json:"created_at"             example:"2025-12-23T00:00:00Z"`
	OtherUser   *UserResponse    `json:"other_user,omitempty"`
	LastMessage *LastMessageInfo `json:"last_message,omitempty"`
	Type        string           `json:"type"                   example:"private"`
	ID          int64            `json:"id"                     example:"123"`
}

func (dto *ChatResponse) Fill(chat model.Chat, otherUser model.User) {
	dto.ID = chat.ID
	dto.Type = string(chat.Type)
	dto.CreatedAt = chat.CreatedAt

	if chat.Type == model.ChatTypePrivate && !otherUser.IsEmpty() {
		other := &UserResponse{}
		other.Fill(otherUser)
		dto.OtherUser = other
	}
}

type LastMessageInfo struct {
	CreatedAt time.Time    `json:"created_at" example:"2025-12-23T00:15:30Z"`
	Text      string       `json:"text"       example:"Привет!"`
	Sender    UserResponse `json:"sender"`
}

func (dto *LastMessageInfo) FillFromModel(msg model.LastMessage) {
	dto.Text = msg.Text
	dto.Sender.Fill(msg.Sender)
	dto.CreatedAt = msg.CreatedAt
}

type ChatsListResponse struct {
	Chats  []ChatResponse `json:"chats"`
	Total  int            `json:"total"  example:"42"`
	Limit  int            `json:"limit"  example:"20"`
	Offset int            `json:"offset" example:"0"`
}

func (dto *ChatsListResponse) Fill(chats []model.ChatWithDetails, total, limit, offset int) {
	dto.Chats = make([]ChatResponse, len(chats))
	for i, chat := range chats {
		dto.Chats[i].FillFromDetails(chat)
	}

	dto.Total = total
	dto.Limit = limit
	dto.Offset = offset
}

func (dto *ChatResponse) FillFromDetails(chat model.ChatWithDetails) {
	dto.ID = chat.ID
	dto.Type = string(chat.Type)
	dto.CreatedAt = chat.CreatedAt

	if chat.Type == model.ChatTypePrivate && !chat.OtherUser.IsEmpty() {
		other := &UserResponse{}
		other.Fill(chat.OtherUser)
		dto.OtherUser = other
	}

	if !chat.LastMessage.IsEmpty() {
		dto.LastMessage = &LastMessageInfo{}
		dto.LastMessage.FillFromModel(chat.LastMessage)
	}
}
