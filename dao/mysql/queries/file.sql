# -- name: CreateFile :one
# START TRANSACTION;
# INSERT INTO files (
#     file_name, file_type, file_size, `key`, url, relation_id, account_id
# ) VALUES (
#              ?, ?, ?, ?, ?, ?, ?
#          );
# SELECT * FROM files
# WHERE file_id = LAST_INSERT_ID();
# COMMIT;


-- name: CreateFile :exec
insert into files (file_name, file_type, file_size, `key`, url, relation_id, account_id)
values(?,?,?,?,?,?,?);

-- name: GetCreateFile :one
select *
from files
where `key`=?;

-- name: DeleteFileByID :exec
delete
from files
where id = ?;

-- name: GetFileKeyByID :one
select `key`
from files
where id = ?;

-- name: GetFileByRelation :many
select *
from files
where relation_id = ?;

-- name: GetFileDetailsByID :one
select *
from files
where id = ?;

-- name: GetGroupAvatar :one
SELECT *
FROM files
WHERE relation_id = ?
  AND account_id IS NULL;

-- name: UpdateGroupAvatar :exec
UPDATE files
SET url = ?
WHERE relation_id = ? AND file_name = 'groupAvatar';

-- name: GetFileByRelationID :many
select id, file_name, file_type, file_size, `key`, url, relation_id, account_id, create_at
from files
where relation_id = ?;