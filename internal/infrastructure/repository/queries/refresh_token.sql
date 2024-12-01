-- refresh_token.sql


-- name: GetRefreshTokenByUserID :one
SELECT
    sqlc.embed(ref_t)
FROM 
    refresh_token ref_t
WHERE
    ref_t.user_id = $1
    AND ref_t.deleted_at IS NULL
    AND ref_t.expired_at > NOW();


-- name: ListRefreshToken :many
SELECT
    sqlc.embed(ref_t)
FROM 
    refresh_token ref_t
WHERE
    (sqlc.narg('user_ids')::UUID[] IS NULL OR u.id = ANY(sqlc.narg('user_ids')::UUID[]))
    AND ref_t.deleted_at IS NULL
    AND ref_t.expired_at > NOW()
LIMIT sqlc.narg('limit')
OFFSET sqlc.arg('offset');
    

-- name: CreateOrUpdateRefreshToken :exec
WITH updated AS (
    UPDATE refresh_token
    SET deleted_at = NOW()
    WHERE user_id = $1
      AND token != $2
      AND deleted_at IS NULL
    RETURNING *
)
INSERT INTO refresh_token (
    user_id, 
    token, 
    expired_at
)
VALUES ($1, $2, $3)
ON CONFLICT (token) DO NOTHING;


-- name: DeleteRefreshTokenByUserID :exec
UPDATE 
    refresh_token
SET 
    deleted_at = COALESCE(deleted_at, NOW())
WHERE 
    user_id = $1;