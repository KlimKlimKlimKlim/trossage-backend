SELECT
    COUNT(*) AS total
FROM
    chats AS c
    INNER JOIN chat_participants AS cp ON c.id = cp.chat_id
WHERE
    cp.user_id = $1
    AND c.deleted_at IS NULL;

