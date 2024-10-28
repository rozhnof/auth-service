package postgres_user_queries

const GetByID = `
	SELECT     
		id, 
		username,
		hash_password
	FROM 
		users
	WHERE
		id = $1
`
