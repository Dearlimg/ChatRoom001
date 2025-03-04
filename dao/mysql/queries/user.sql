-- name: CreateUser :exec
INSERT INTO users (email, password)
VALUES (?, ?);
# -- 这里假设使用编程语言调用，通过执行获取自增 ID 的操作，若在 MySQL 客户端可直接执行此语句
# SELECT LAST_INSERT_ID() as id, email, password, create_at
# FROM users
# WHERE id = LAST_INSERT_ID();

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