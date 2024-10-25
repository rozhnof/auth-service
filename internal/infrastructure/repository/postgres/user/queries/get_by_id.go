package postgres_user_queries

const GetByIDQuery = `
	SELECT     
		id, 
		email
	FROM 
		users
	WHERE
		id = $1
`
