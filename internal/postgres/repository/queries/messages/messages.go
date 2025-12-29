package messages

import (
	_ "embed"
)

//go:embed create_message.sql
var CreateMessageQuery string
