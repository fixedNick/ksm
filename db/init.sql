CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE chat_type_id AS ENUM ('user', 'public', 'private', 'public_channel', 'private_channel');

CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(48) NOT NULL UNIQUE,
    password_hash VARCHAR(60) NOT NULL, -- BCRYPT
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE chats (
    chat_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_owner UUID REFERENCES users(user_id) ON DELETE SET NULL,
    chat_name VARCHAR(48) NOT NULL DEFAULT 'Новый чат',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    chat_type chat_type_id NOT NULL DEFAULT 'user'
);

CREATE TABLE chat_members (
    chat_id UUID REFERENCES chats(chat_id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE messages (
    message_id BIGSERIAL PRIMARY KEY,
    chat_id UUID REFERENCES chats(chat_id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(user_id) ON DELETE SET NULL,
    
    content TEXT DEFAULT '',
    attachments JSONB DEFAULT '[]'::jsonb, -- Files, voices, imgs, etc.
    reactions JSONB DEFAULT '{}'::jsonb,   -- Reactions on message only
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE message_reads (
    message_id BIGINT REFERENCES messages(message_id) ON DELETE CASCADE, 
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    read_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (message_id, user_id)
);

CREATE TABLE device_tokens (
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    updated_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, token)
);

CREATE INDEX idx_messages_chat_id_created_at ON messages(chat_id, created_at DESC);
CREATE INDEX idx_message_reads_user_id ON message_reads(user_id);
CREATE INDEX idx_chat_members_user_id_include_chat ON chat_members(user_id) INCLUDE (chat_id);