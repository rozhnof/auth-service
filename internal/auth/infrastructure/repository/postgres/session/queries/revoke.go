package postgres_session_queries

const Revoke = `
	UPDATE 
		session 
	SET  
		is_revoked = TRUE
	WHERE 
		user_id = $1
`
