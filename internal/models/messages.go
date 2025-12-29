package models

import (
	"time"
)

type Message struct {
	CreatedAt time.Time
	Text      string
	ID        int64
	ChatID    int64
	SenderID  int64
}
