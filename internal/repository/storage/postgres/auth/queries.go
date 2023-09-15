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
	getByIDQuery         = `SELECT * FROM users WHERE id = $1`
	getUsersQuery        = `SELECT * FROM users`
	updateLastLoginQuery = `UPDATE users SET last_login_at = now() WHERE id = $1`
	updateUserQuery      = `SET
    first_name = $1,
    last_name = $2,
    email = $3,
    password = $4,
    role = COALESCE(NULLIF($5, ''), 'user'),
    profile_picture_url = $6,
    phone = $7,
    city = $8,
    gender = $9,
    birthday = $10,
    age = $11,
    updated_at = now()
  WHERE
    id = $12
  RETURNING *`
	deleteUserQuery = `DELETE FROM users WHERE id = $1`
)
