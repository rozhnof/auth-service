CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL,
    hash_password TEXT NOT NULL,
    deleted_at TIMESTAMP
);

CREATE TABLE session (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    refresh_token VARCHAR(255),
    expired_at TIMESTAMP,
    is_revoked BOOLEAN,
    deleted_at TIMESTAMP
);

