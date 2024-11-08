package postgres_user_queries

const List = `
	SELECT     
		id, 
		username,
		hash_password
	FROM 
		users
`
