CREATE TABLE users (
    id bigserial PRIMARY KEY,
    login VARCHAR(20) NOT NULL,
    password_hash text NOT NULL,
    display_name varchar(20) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_users_login ON users (login);

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

