INSERT INTO chats (chat_type)
    VALUES ($1)
RETURNING
    id, chat_type, created_at, updated_at;

