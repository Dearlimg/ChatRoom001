-- name: CreateAccount :exec
insert into accounts (id, user_id, name, avatar, gender, signature)
values (?, ?, ?, ?, ?, ?);

-- name: DeleteAccount :exec
delete
from accounts
where id = ?;

-- name: DeleteAccountByUserID :exec
delete
from accounts
where user_id = ?;

-- name: UpdateAccount :exec
update accounts
set name = ?,
    gender=?,
    signature=?
where id =?;


-- name: UpdateAccountAvatar :exec
update accounts
set avatar = ?
where id = ?;

-- name: GetAccountByID :many
SELECT
    a.*,
    r.id AS relation_id
FROM (
         SELECT *
         FROM accounts
         WHERE accounts.user_id = ?
     ) a
         LEFT JOIN relations r
                   ON r.relation_type = 'friend'
                       AND (
                          (r.account1_id = a.id AND r.account2_id = ?) -- 当前用户ID
                              OR
                          (r.account1_id = ? AND r.account2_id = a.id) -- 当前用户ID
                          )
LIMIT 1;


-- name: GetAccountByUserID :many
select id,name,avatar,gender
from accounts
where user_id = ?;

-- name: ExistAccountByID :one
select exists(
    select 1
    from accounts
    where id =?
);


-- name: CountAccountByUserID :one
SELECT COUNT(id) AS count
FROM accounts
WHERE user_id = ?;

-- name: ExistsAccountByNameAndUserID :one
select exists(
    select 1
    from accounts
    where user_id=?
    and name =?
);

-- name: GetAccountsByName :many
SELECT
    a.*,
    r.id AS relation_id,
    (SELECT COUNT(*) FROM accounts WHERE name LIKE CONCAT('%', ?, '%')) AS total
FROM (
         SELECT id, name, avatar, gender
         FROM accounts
         WHERE name LIKE CONCAT('%', ?, '%')
     ) AS a
         LEFT JOIN relations r
                   ON r.relation_type = 'friend'
                       AND (
                          (r.account1_id = a.id AND r.account2_id = ?)
                              OR
                          (r.account1_id = ? AND r.account2_id = a.id)
                          )
LIMIT ? OFFSET ?;
