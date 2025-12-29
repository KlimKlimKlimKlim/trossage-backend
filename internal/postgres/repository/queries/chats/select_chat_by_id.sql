SELECT
    id,
    chat_type,
    created_at,
    updated_at,
    COALESCE(deleted_at, '0001-01-01'::timestamptz) AS deleted_at
FROM
    chats
WHERE
    id = $1
    AND deleted_at IS NULL;

