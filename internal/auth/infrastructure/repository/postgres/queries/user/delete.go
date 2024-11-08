package postgres_user_queries

const Delete = `
UPDATE 
    users
SET 
    deleted_at = COALESCE(deleted_at, NOW())
WHERE 
    id = $1
RETURNING 
    deleted_at;
`
