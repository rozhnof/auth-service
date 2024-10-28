package postgres_session_queries

const UpdateQuery = `
	UPDATE 
		session 
	SET  
		user_id = $2,
		refresh_token = $3,
		expired_at = $4,
		is_revoked = $5
	WHERE 
		id = $1
	RETURNING 
		id,
		user_id,
		refresh_token,
		expired_at,
		is_revoked
`
