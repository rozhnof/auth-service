package postgres_session_queries

const Create = `
	INSERT INTO session (
		user_id,
		refresh_token,
		expired_at
	) VALUES (
	 	$1, $2, $3
	)
	RETURNING 
		id,
		user_id,
		refresh_token,
		expired_at,
		is_revoked
`
