-- squawk-ignore-file ban-drop-table, require-concurrent-index-deletion
SET lock_timeout = '2s';

SET statement_timeout = '5s';

DROP TRIGGER IF EXISTS set_updated_at_users ON users;

DROP FUNCTION IF EXISTS update_updated_at_users ();

DROP INDEX IF EXISTS idx_users_login;

DROP INDEX IF EXISTS idx_users_login_prefix;

DROP TABLE IF EXISTS users;

