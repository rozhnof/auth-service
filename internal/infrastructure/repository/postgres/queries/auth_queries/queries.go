package auth_queries

const CreateQuery = `
	INSERT INTO users (
		id, 
		email,
		token_hash
	) VALUES (
	 	$1, $2, $3
	)
	RETURNING id
`

const GetByIDQuery = `
	SELECT     
		id, 
		email,
		token_hash
	FROM 
		users
	WHERE
		id = $1
`

const GetByEmailQuery = `
	SELECT     
		id, 
		email,
		token_hash
	FROM 
		users
	WHERE
		email = $1
`

const ListQuery = `
	SELECT     
		id, 
		email,
		token_hash
	FROM 
		users
`

const UpdateQuery = `
	UPDATE 
		users 
	SET  
		email = $2,
		token_hash = $3
	WHERE 
		id = $1
`

const DeleteQuery = `
	DELETE 
		FROM user 
	WHERE 
		id = $1
`
