package postgres_user_queries

const GetByUsernameQuery = `
	SELECT     
		id, 
		username
	FROM 
		users
	WHERE
		username = $1
`
