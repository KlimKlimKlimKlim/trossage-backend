package message

import (
	_ "embed"
)

//go:embed create_message.sql
var CreateMessageQuery string

//go:embed select_messages.sql
var SelectMessagesQuery string

//go:embed count_chat_messages.sql
var CountChatMessagesQuery string
