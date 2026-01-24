SELECT
    user_id
FROM
    chat_participants
WHERE
    chat_id = $1;

