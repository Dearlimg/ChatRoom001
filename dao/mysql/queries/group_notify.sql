-- name: CreateGroupNotify :exec
insert into group_notify
(relation_id, msg_content, msg_expand, account_id, create_at, read_ids)
values (?,?,?,?,?,?);

# -- name: CreateGroupNotifyReturn :one
# select *
# from group_notify
# where create_at=?
# and relation_id=?
# and account_id=?;


-- name: CreateGroupNotifyReturn :one
SELECT *
FROM group_notify
WHERE id = LAST_INSERT_ID();

-- name: UpdateGroupNotify :exec
update group_notify
set relation_id=?,
    msg_content=?,
    msg_expand=?,
    account_id=?,
    create_at=?,
    read_ids=?
where id=?;

-- name: UpdateGroupNotifyReturn :one
select *
from group_notify
where id=?;

-- name: GetGroupNotifyByID :many
select id, relation_id, msg_content, msg_expand, account_id, create_at, read_ids
from group_notify
where id = ?;

-- name: DeleteGroupNotify :exec
delete
from group_notify
where id=?;

