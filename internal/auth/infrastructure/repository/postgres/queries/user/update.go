package postgres_user_queries

const Update = `
	UPDATE 
		users 
	SET  
		username = $2,
		hash_password = $3
	WHERE 
		id = $1
	RETURNING 
		id,
		username,
		hash_password
`
