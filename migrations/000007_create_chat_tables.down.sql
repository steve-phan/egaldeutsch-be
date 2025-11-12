-- Drop indexes
DROP INDEX IF EXISTS idx_chat_messages_created_at;
DROP INDEX IF EXISTS idx_chat_messages_user_id;
DROP INDEX IF EXISTS idx_chat_messages_room_id;
DROP INDEX IF EXISTS idx_chat_rooms_active;
DROP INDEX IF EXISTS idx_chat_rooms_created_by;

-- Drop tables
DROP TABLE IF EXISTS chat_messages;
DROP TABLE IF EXISTS chat_rooms;
