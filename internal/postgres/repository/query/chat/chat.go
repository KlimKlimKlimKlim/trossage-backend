package chat

import (
	_ "embed"
)

//go:embed select_chat_between_users.sql
var SelectChatBetweenUsersQuery string

//go:embed insert_chat.sql
var InsertChatQuery string

//go:embed insert_chat_participants.sql
var InsertChatParticipantsQuery string

//go:embed select_user_chats.sql
var SelectUserChatsQuery string

//go:embed count_user_chats.sql
var CountUserChatsQuery string

//go:embed is_user_member.sql
var IsUserMemberQuery string

//go:embed select_chat_by_id.sql
var SelectChatByIDQuery string

//go:embed select_chat_members.sql
var SelectChatMembersQuery string
