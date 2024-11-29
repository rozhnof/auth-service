-- queries.sql


-- name: GetUserByID :one
SELECT 
    sqlc.embed(u),
    rt.id AS refresh_token_id,
    rt.refresh_token AS refresh_token,
    rt.expired_at AS refresh_token_expired_at
FROM 
    users u
LEFT JOIN 
    refresh_token rt ON u.id = rt.user_id
WHERE 
    u.id = $1
    AND u.deleted_at IS NULL
    AND rt.deleted_at IS NULL;


-- name: GetUserByEmail :one
SELECT 
    sqlc.embed(u),
    rt.id AS refresh_token_id,
    rt.refresh_token AS refresh_token,
    rt.expired_at AS refresh_token_expired_at
FROM 
    users u
LEFT JOIN 
    refresh_token rt ON u.id = rt.user_id
WHERE 
    u.email = $1
    AND u.deleted_at IS NULL
    AND rt.deleted_at IS NULL;


-- name: GetUserByRefreshToken :one
SELECT 
    sqlc.embed(u),
    rt.id AS refresh_token_id,
    rt.refresh_token AS refresh_token,
    rt.expired_at AS refresh_token_expired_at
FROM 
    users u
LEFT JOIN 
    refresh_token rt ON u.id = rt.user_id
WHERE 
    rt.refresh_token = $1
    AND u.deleted_at IS NULL
    AND rt.deleted_at IS NULL;


-- name: ListUser :many
SELECT 
    sqlc.embed(u),
    rt.id AS refresh_token_id,
    rt.refresh_token AS refresh_token,
    rt.expired_at AS refresh_token_expired_at
FROM 
    users u
LEFT JOIN 
    refresh_token rt ON u.id = rt.user_id
WHERE 
    (sqlc.narg('user_ids')::UUID[] IS NULL OR u.id = ANY(sqlc.narg('user_ids')::UUID[]))
    AND u.deleted_at IS NULL
    AND rt.deleted_at IS NULL
LIMIT sqlc.narg('limit')
OFFSET sqlc.arg('offset');


-- name: CreateUser :exec
INSERT INTO users (
    id,
    email,
    hash_password
) VALUES (
    $1, $2, $3
);


-- name: UpdateUser :exec
UPDATE 
    users
SET  
    email = $2,
    hash_password = $3
WHERE 
    users.id = $1;


-- name: DeleteUser :exec
UPDATE 
    users
SET 
    deleted_at = COALESCE(deleted_at, NOW())
WHERE 
    id = $1;


-- name: CreateOrUpdateRefreshToken :exec
WITH updated AS (
    UPDATE refresh_token
    SET deleted_at = NOW()
    WHERE user_id = $1
      AND refresh_token != $2
      AND deleted_at IS NULL
    RETURNING *
)
INSERT INTO refresh_token (
    user_id, 
    refresh_token, 
    expired_at
)
VALUES ($1, $2, $3)
ON CONFLICT (refresh_token) DO NOTHING;


-- name: DeleteRefreshTokenByUserID :exec
UPDATE 
    refresh_token
SET 
    deleted_at = COALESCE(deleted_at, NOW())
WHERE 
    user_id = $1;


-- name: CreateOutboxMessage :exec
INSERT INTO outbox (
    key,
    value,
    topic
) VALUES (
    $1, $2, $3
);

-- name: ReadOutboxMessageList :many
WITH messages AS (
    SELECT 
        id, 
        key, 
        value, 
        topic
    FROM 
        outbox o
    WHERE 
        o.topic = $1
        AND o.is_read = FALSE
        AND o.deleted_at IS NULL
    ORDER BY o.created_at
    LIMIT $2
)
UPDATE 
    outbox
SET 
    is_read = TRUE
FROM 
    messages m
WHERE 
    outbox.id = m.id
RETURNING 
    m.key, 
    m.value, 
    m.topic;
