-- squawk-ignore-file ban-drop-table, require-concurrent-index-deletion
SET lock_timeout = '2s';

SET statement_timeout = '5s';

DROP TRIGGER IF EXISTS update_chat_on_message ON messages;

DROP FUNCTION IF EXISTS update_chat_last_message ();

DROP INDEX IF EXISTS idx_messages_chat_created;

DROP TABLE IF EXISTS messages;

DROP INDEX IF EXISTS idx_chat_participants_user;

DROP TABLE IF EXISTS chat_participants;

DROP INDEX IF EXISTS idx_chats_last_message;

DROP TRIGGER IF EXISTS set_updated_at_chats ON chats;

DROP FUNCTION IF EXISTS update_updated_at_chats ();

DROP TABLE IF EXISTS chats;

