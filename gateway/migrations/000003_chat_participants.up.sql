CREATE TABLE chat_participants (
                                   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                   chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
                                   user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                   joined_at TIMESTAMPTZ DEFAULT now(),
                                   UNIQUE (chat_id, user_id)
);
