-- name: CreateUser :one
INSERT INTO users (
  email,
  password
)	VALUES (
	$1, $2
) RETURNING *;


-- name: SelectUser :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;


-- name: DeleteUser :exec
DELETE FROM users WHERE email = $1;