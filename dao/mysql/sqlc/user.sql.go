// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO users (email, password,create_at)
VALUES (?, ?,now())
`

type CreateUserParams struct {
	Email    string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg *CreateUserParams) error {
	_, err := q.exec(ctx, q.createUserStmt, createUser, arg.Email, arg.Password)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteUserStmt, deleteUser, id)
	return err
}

const existEmail = `-- name: ExistEmail :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = ?) as result
`

func (q *Queries) ExistEmail(ctx context.Context, email string) (bool, error) {
	row := q.queryRow(ctx, q.existEmailStmt, existEmail, email)
	var result bool
	err := row.Scan(&result)
	return result, err
}

const existsUserByID = `-- name: ExistsUserByID :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = ?) as result
`

func (q *Queries) ExistsUserByID(ctx context.Context, id int64) (bool, error) {
	row := q.queryRow(ctx, q.existsUserByIDStmt, existsUserByID, id)
	var result bool
	err := row.Scan(&result)
	return result, err
}

const getAcountIDsByUserID = `-- name: GetAcountIDsByUserID :many
SELECT id
FROM accounts
WHERE user_id = ?
`

func (q *Queries) GetAcountIDsByUserID(ctx context.Context, userID int64) ([]int64, error) {
	rows, err := q.query(ctx, q.getAcountIDsByUserIDStmt, getAcountIDsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int64{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllEmail = `-- name: GetAllEmail :many
SELECT email
FROM users
`

func (q *Queries) GetAllEmail(ctx context.Context) ([]string, error) {
	rows, err := q.query(ctx, q.getAllEmailStmt, getAllEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		items = append(items, email)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, password, create_at
FROM users
WHERE email = ?
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := q.queryRow(ctx, q.getUserByEmailStmt, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreateAt,
	)
	return &i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, email, password, create_at
FROM users
WHERE id = ?
LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (*User, error) {
	row := q.queryRow(ctx, q.getUserByIDStmt, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreateAt,
	)
	return &i, err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET email = ?,
    password = ?
WHERE id = ?
`

type UpdateUserParams struct {
	Email    string
	Password string
	ID       int64
}

func (q *Queries) UpdateUser(ctx context.Context, arg *UpdateUserParams) error {
	_, err := q.exec(ctx, q.updateUserStmt, updateUser, arg.Email, arg.Password, arg.ID)
	return err
}
