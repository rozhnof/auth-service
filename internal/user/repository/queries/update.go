package queries

const UpdateQuery = `
	UPDATE 
		users 
	SET  
		username = $2
	WHERE 
		id = $1
	RETURNING 
		id,
		username
`
