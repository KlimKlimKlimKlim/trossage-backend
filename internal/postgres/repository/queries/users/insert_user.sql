INSERT INTO users (login, password_hash, display_name)
    VALUES ($1, $2, $3)
RETURNING
    id, created_at, updated_at;

