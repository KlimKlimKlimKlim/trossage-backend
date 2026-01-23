package hub

import (
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/dto"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/model"
)

type EventType string

const (
	EventTypeTyping      EventType = "typing"
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

func NewTypingEvent(senderID, chatID int64, operations []dto.TypingOperation) *Event {
	var typingDTO dto.TypingEventResponse
	typingDTO.Fill(senderID, chatID, operations)

	return newEvent(EventTypeTyping, typingDTO)
}

func NewNewMessageEvent(message model.Message) *Event {
	var messageDTO dto.MessageResponse
	messageDTO.Fill(message)

	return newEvent(EventTypeNewMessage, messageDTO)
}

func NewChatCreatedEvent(chat model.Chat, user model.User) *Event {
	var chatDTO dto.ChatResponse
	chatDTO.Fill(chat, user)

	return newEvent(EventTypeChatCreated, chatDTO)
}

func NewUserLogoutEvent() *Event {
	return newEvent(EventTypeUserLogout, nil)
}
