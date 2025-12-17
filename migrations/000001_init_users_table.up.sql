CREATE TABLE users (
    id bigserial PRIMARY KEY,
    login VARCHAR(30) NOT NULL,
    password_hash text NOT NULL,
    display_name varchar(30) NOT NULL,
    deleted_at timestamptz DEFAULT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_users_login ON users (login)
WHERE
    deleted_at IS NULL;

CREATE OR REPLACE FUNCTION update_updated_at_users ()
    RETURNS TRIGGER
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_users
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_users ();

