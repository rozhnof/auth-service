-- outbox.sql


-- name: CreateOutboxMessage :exec
INSERT INTO outbox (
    key,
    value,
    topic
) VALUES (
    $1, $2, $3
);

-- name: CreateOutboxMessages :copyfrom
INSERT INTO outbox (
    key,
    value,
    topic
) VALUES (
    $1, $2, $3
);

-- name: ReadOutboxMessages :many
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