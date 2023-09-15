package auth

const (
	createUserQuery = `
  INSERT INTO users (
    first_name,
    last_name,
    email,
    password,
    role,
    profile_picture_url,
    phone,
    city,
    gender,
    birthday,
    age,
    created_at,
    updated_at,
    last_login_at
  )
  VALUES (
    $1, $2, $3, $4, COALESCE(NULLIF($5, ''), 'user'), $6, $7, $8, $9, $10, $11, now(), now(), now()
  )
  RETURNING *
`
	getByEmailQuery      = `SELECT * FROM users WHERE email = $1`
	getByPhoneQuery      = `SELECT * FROM users WHERE phone = $1`
	getUsersQuery        = "SELECT * FROM users"
	updateLastLoginQuery = "UPDATE users SET last_login_at = now() WHERE id = $1"
)
