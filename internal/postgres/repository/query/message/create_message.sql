INSERT INTO messages (chat_id, sender_id, message_text)
    VALUES ($1, $2, $3)
RETURNING
    id, chat_id, sender_id, message_text, created_at;

