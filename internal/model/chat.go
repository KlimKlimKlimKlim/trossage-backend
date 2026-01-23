package model

import (
	"time"
)

type ChatType string

const (
	ChatTypePrivate ChatType = "private"
)

type Chat struct {
	DeletedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Type      ChatType
	Title     string
	ID        int64
}

type ChatWithDetails struct {
	LastMessage LastMessage
	Chat
	OtherUser User
}

type LastMessage struct {
	CreatedAt time.Time
	Text      string
	Sender    User
}

func (m LastMessage) IsEmpty() bool {
	return m.Sender.ID == 0
}
