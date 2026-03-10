CREATE TABLE chats (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       is_private BOOLEAN NOT NULL,           -- Признак: личный чат или групповой
                       name VARCHAR(255),                     -- Название чата (для групповых)
                       created_at TIMESTAMPTZ DEFAULT now()   -- Дата создания
);
