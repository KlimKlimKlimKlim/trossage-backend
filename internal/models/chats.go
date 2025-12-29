package models

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
