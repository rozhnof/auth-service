package postgres_user_queries

const Create = `
	INSERT INTO users (
		username,
		hash_password
	) VALUES (
	 	$1, $2
	)
	RETURNING 
		id,
		username,
		hash_password
`
