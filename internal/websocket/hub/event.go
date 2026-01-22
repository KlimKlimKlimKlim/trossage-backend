package hub

import "encoding/json"

type EventType string

const (
	EventTypeNewMessage     EventType = "new_message"
	EventTypeMessageEdited  EventType = "message_edited"
	EventTypeMessageDeleted EventType = "message_deleted"
	EventTypeChatCreated    EventType = "chat_created"
	EventTypeChatUpdated    EventType = "chat_updated"
)

// ServerEvent - событие от сервера к клиенту.
type ServerEvent struct {
	Type EventType       `json:"type"`
	Data json.RawMessage `json:"data"`
}

// NewMessageEvent - данные события нового сообщения.
type NewMessageEvent struct {
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	ChatID    int64  `json:"chat_id"`
	MessageID int64  `json:"message_id"`
	SenderID  int64  `json:"sender_id"`
}

// MessageDeletedEvent - сообщение удалено.
type MessageDeletedEvent struct {
	ChatID    int64 `json:"chat_id"`
	MessageID int64 `json:"message_id"`
}

// ChatCreatedEvent - создан новый чат.
type ChatCreatedEvent struct {
	ChatID int64 `json:"chat_id"`
	UserID int64 `json:"user_id"` // С кем создан чат.
}

// NewServerEvent создает событие для отправки клиенту.
func NewServerEvent(eventType EventType, data any) (*ServerEvent, error) {
	rawData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &ServerEvent{
		Type: eventType,
		Data: rawData,
	}, nil
}
