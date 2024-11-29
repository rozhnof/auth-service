CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(50) NOT NULL UNIQUE,
    hash_password TEXT NOT NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE refresh_token (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID REFERENCES users (id) NOT NULL,
    refresh_token VARCHAR(255) NOT NULL UNIQUE,
    expired_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    key UUID NOT NULL UNIQUE,
    value JSON NOT NULL,
    topic TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);


CREATE UNIQUE INDEX idx_users_email ON users (email); 
CREATE INDEX idx_refresh_token_user_id ON refresh_token (user_id); 