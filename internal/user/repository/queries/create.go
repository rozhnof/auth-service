package queries

const CreateQuery = `
	INSERT INTO users (
		username,
		password
	) VALUES (
	 	$1
	)
	RETURNING 
		id,
		username
`
