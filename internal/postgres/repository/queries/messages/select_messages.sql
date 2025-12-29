SELECT
    id,
    chat_id,
    sender_id,
    message_text,
    created_at
FROM
    messages
WHERE
    chat_id = $1
ORDER BY
    created_at DESC
LIMIT $2 OFFSET $3;

