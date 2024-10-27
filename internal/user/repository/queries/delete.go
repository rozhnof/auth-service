package queries

const DeleteQuery = `
UPDATE 
    users
SET 
    deleted_at = COALESCE(deleted_at, NOW())
WHERE 
    id = $1
RETURNING 
    deleted_at;
`
