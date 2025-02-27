-- name: CreateAccount :exec
insert into accounts (id, user_id, name, avatar, gender, signature)
values (?, ?, ?, ?, ?, ?);

