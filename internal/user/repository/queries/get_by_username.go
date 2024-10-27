package queries

const GetByUsernameQuery = `
	SELECT     
		id, 
		username
	FROM 
		users
	WHERE
		username = $1
`
