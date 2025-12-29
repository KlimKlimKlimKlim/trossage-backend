package chats

import (
	_ "embed"
)

//go:embed select_chat_between_users.sql
var SelectChatBetweenUsersQuery string

//go:embed insert_chat.sql
var InsertChatQuery string

//go:embed insert_chat_participants.sql
var InsertChatParticipantsQuery string
