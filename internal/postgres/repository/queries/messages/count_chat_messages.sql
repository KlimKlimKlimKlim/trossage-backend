SELECT
    COUNT(*) AS total
FROM
    messages
WHERE
    chat_id = $1;

