-- name: CreateUser :exec
INSERT INTO users (email, password)
VALUES (?, ?);

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = ?;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?
LIMIT 1;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = ?
LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
SET email = ?,
    password = ?
WHERE id = ?;

-- name: ExistEmail :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = ?) as result;

-- name: GetAllEmail :many
SELECT email
FROM users;

-- name: ExistsUserByID :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = ?) as result;

-- name: GetAcountIDsByUserID :many
SELECT id
FROM accounts
WHERE user_id = ?;