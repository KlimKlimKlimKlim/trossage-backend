SELECT
    EXISTS (
        SELECT
            1 AS one
        FROM
            chat_participants AS cp
            INNER JOIN chats AS c ON cp.chat_id = c.id
        WHERE
            cp.chat_id = $1
            AND cp.user_id = $2
            AND c.deleted_at IS NULL) AS is_member;

