SELECT
    c.id
FROM
    chats AS c
    INNER JOIN chat_participants AS cp1 ON c.id = cp1.chat_id
    INNER JOIN chat_participants AS cp2 ON c.id = cp2.chat_id
WHERE
    cp1.user_id = $1
    AND cp2.user_id = $2
    AND c.deleted_at IS NULL
ORDER BY
    c.id
LIMIT 1;

