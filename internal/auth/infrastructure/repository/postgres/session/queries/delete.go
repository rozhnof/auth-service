package postgres_session_queries

const Delete = `
UPDATE 
    session
SET 
    deleted_at = COALESCE(deleted_at, NOW())
WHERE 
    id = $1
RETURNING 
    deleted_at;
`
