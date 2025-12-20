-- squawk-ignore-file ban-drop-table, require-concurrent-index-deletion
SET lock_timeout = '2s';

SET statement_timeout = '5s';

DROP INDEX IF EXISTS idx_refresh_tokens_revoked;

DROP INDEX IF EXISTS idx_refresh_tokens_expired;

DROP INDEX IF EXISTS idx_refresh_tokens_user_active;

DROP INDEX IF EXISTS idx_refresh_tokens_user_token;

DROP TABLE IF EXISTS refresh_tokens;

