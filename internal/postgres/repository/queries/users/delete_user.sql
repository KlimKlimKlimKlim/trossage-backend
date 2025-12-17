UPDATE
    users
SET
    login = CONCAT('deleted_user_', id),
    display_name = 'Deleted User',
    password_hash = '',
    deleted_at = NOW()
WHERE
    id = $1
    AND deleted_at IS NULL;

