SELECT
    c.id,
    c.chat_type,
    c.created_at,
    c.last_message_text,
    c.last_message_sender_id,
    c.last_message_at,
    last_sender.user_login AS last_sender_login,
    last_sender.display_name AS last_sender_display_name,
    last_sender.created_at AS last_sender_created_at,
    other_user.id AS other_user_id,
    other_user.user_login AS other_user_login,
    other_user.display_name AS other_user_display_name,
    other_user.created_at AS other_user_created_at
FROM
    chats AS c
    INNER JOIN chat_participants AS cp ON c.id = cp.chat_id
    LEFT JOIN chat_participants AS other_cp ON c.id = other_cp.chat_id
        AND other_cp.user_id != $1
    LEFT JOIN users AS other_user ON other_cp.user_id = other_user.id
    LEFT JOIN users AS last_sender ON c.last_message_sender_id = last_sender.id
WHERE
    cp.user_id = $1
    AND c.deleted_at IS NULL
ORDER BY
    c.last_message_at DESC NULLS LAST,
    c.created_at DESC
LIMIT $2 OFFSET $3;

