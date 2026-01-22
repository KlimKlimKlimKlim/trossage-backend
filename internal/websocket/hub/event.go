package hub

import (
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/dto"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
)

type EventType string

const (
	EventTypeNewMessage  EventType = "new_message"
	EventTypeChatCreated EventType = "chat_created"
	EventTypeUserLogout  EventType = "user_logout"
)

type Event struct {
	Data any       `json:"data"`
	Type EventType `json:"type"`
}

func newEvent(eventType EventType, data any) *Event {
	return &Event{
		Type: eventType,
		Data: data,
	}
}

func NewNewMessageEvent(message models.Message) *Event {
	var messageDTO dto.MessageResponse
	messageDTO.Fill(message)

	return newEvent(EventTypeNewMessage, messageDTO)
}

func NewChatCreatedEvent(chat models.Chat, user models.User) *Event {
	var chatDTO dto.ChatResponse
	chatDTO.Fill(chat, user)

	return newEvent(EventTypeChatCreated, chatDTO)
}

func NewUserLogoutEvent() *Event {
	return newEvent(EventTypeUserLogout, nil)
}
