UPDATE
    users
SET
    display_name = $2,
    password_hash = $3
WHERE
    id = $1
    AND deleted_at IS NULL
RETURNING
    updated_at;

