CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL,
    password TEXT NOT NULL,
    refresh_token TEXT NOT NULL
);