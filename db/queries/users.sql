-- name: CreateUser :one
INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: LookupUser :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: DeleteAllUsers :exec
TRUNCATE TABLE users;