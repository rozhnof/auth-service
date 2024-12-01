-- user.sql


-- name: GetUserByID :one
SELECT 
    sqlc.embed(u)
FROM 
    users u
WHERE 
    u.id = $1
    AND u.deleted_at IS NULL;


-- name: GetUserByEmail :one
SELECT 
    sqlc.embed(u)
FROM 
    users u
WHERE 
    u.email = $1
    AND u.deleted_at IS NULL;


-- name: GetUserByRefreshToken :one
SELECT 
    sqlc.embed(u)
FROM 
    users u
JOIN 
    refresh_token ref_t ON u.id = ref_t.user_id
WHERE
    ref_t.token = $1
    AND ref_t.deleted_at IS NULL;


-- name: ListUser :many
SELECT 
    sqlc.embed(u)
FROM 
    users u
WHERE 
    (sqlc.narg('user_ids')::UUID[] IS NULL OR u.id = ANY(sqlc.narg('user_ids')::UUID[]))
    AND u.deleted_at IS NULL
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
    hash_password = $3,
    confirmed = $4
WHERE 
    users.id = $1;


-- name: DeleteUser :exec
UPDATE 
    users
SET 
    deleted_at = COALESCE(deleted_at, NOW())
WHERE 
    id = $1;
