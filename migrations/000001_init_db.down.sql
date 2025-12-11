DROP TRIGGER IF EXISTS set_updated_at_users ON users;

DROP FUNCTION IF EXISTS update_updated_at_users ();

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS refresh_tokens CASCADE;

