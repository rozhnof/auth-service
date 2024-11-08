package postgres_session_queries

const GetByRefreshToken = `
	SELECT     
		id, 
		user_id,
		refresh_token,
		expired_at,
		is_revoked
	FROM 
		session
	WHERE
		refresh_token = $1
`
