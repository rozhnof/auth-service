package postgres_session_queries

const GetByID = `
	SELECT     
		id, 
		user_id,
		refresh_token,
		expired_at,
		is_revoked
	FROM 
		session
	WHERE
		id = $1
`
