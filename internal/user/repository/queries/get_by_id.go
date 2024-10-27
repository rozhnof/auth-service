package queries

const GetByIDQuery = `
	SELECT     
		id, 
		email
	FROM 
		users
	WHERE
		id = $1
`
