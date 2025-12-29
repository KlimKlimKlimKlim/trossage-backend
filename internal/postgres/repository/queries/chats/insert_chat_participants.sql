INSERT INTO chat_participants (chat_id, user_id)
SELECT
    $1 AS chat_id,
    unnest($2::bigint[]) AS user_id;

