CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(50) NOT NULL UNIQUE,
    confirmed BOOL NOT NULL DEFAULT FALSE,
    hash_password CHAR(60) NOT NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE refresh_token (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID REFERENCES users (id) NOT NULL,
    token VARCHAR NOT NULL UNIQUE,
    expired_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE register_token (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID REFERENCES users (id) NOT NULL,
    token VARCHAR NOT NULL UNIQUE,
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