package postgres_session_queries

const Create = `
	INSERT INTO session (
		user_id,
		refresh_token,
		expired_at,
	 	is_revoked
	) VALUES (
	 	$1, $2, $3, $4
	)
	RETURNING 
		id,
		user_id,
		refresh_token,
		expired_at,
		is_revoked
`
