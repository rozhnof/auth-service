package postgres_user_queries

const GetByUsername = `
	SELECT     
		id, 
		username,
		hash_password
	FROM 
		users
	WHERE
		username = $1
`
