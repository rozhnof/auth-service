-- register_token.sql


-- name: GetRegisterTokenByUserID :one
SELECT
    sqlc.embed(reg_t)
FROM 
    register_token reg_t
WHERE
    reg_t.user_id = $1
    AND reg_t.deleted_at IS NULL
    AND reg_t.expired_at > NOW();


-- name: ListRegisterToken :many
SELECT
    sqlc.embed(reg_t)
FROM 
    register_token reg_t
WHERE
    (sqlc.narg('user_ids')::UUID[] IS NULL OR u.id = ANY(sqlc.narg('user_ids')::UUID[]))
    AND reg_t.deleted_at IS NULL
    AND reg_t.expired_at > NOW()
LIMIT sqlc.narg('limit')
OFFSET sqlc.arg('offset');


-- name: CreateOrUpdateRegisterToken :exec
WITH updated AS (
    UPDATE register_token
    SET deleted_at = NOW()
    WHERE user_id = $1
      AND token != $2
      AND deleted_at IS NULL
    RETURNING *
)
INSERT INTO register_token (
    user_id, 
    token, 
    expired_at
)
VALUES ($1, $2, $3)
ON CONFLICT (token) DO NOTHING;


-- name: DeleteRegisterTokenByUserID :exec
UPDATE 
    register_token
SET 
    deleted_at = COALESCE(deleted_at, NOW())
WHERE 
    user_id = $1;

