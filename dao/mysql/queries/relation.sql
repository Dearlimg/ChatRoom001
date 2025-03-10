-- name: CreateGroupRelation :one
INSERT INTO relations (relation_type, group_name, group_description, group_avatar)
value ('group',?,?,?);

-- name: DeleteFriendRelation :one
insert into relations(relation_type,account1_id,account2_id)
value ('friend',?,?);

-- name: DeleteRelation :exec
delete
from relations
where id=?;

-- name: DeleteFriendRelationByAccountID :many
delete
from relations
where relation_type ='friend'
and (account1_id=?);

-- 第一步：查询待删除的ID
SELECT id
FROM relations
WHERE
    relation_type = 'friend'
  AND (
    account1_id = :account2_id
        OR account2_id = :account1_id
    );

-- 第二步：执行删除操作
DELETE FROM relations
WHERE
    relation_type = 'friend'
  AND (
    account1_id = :account1_id
        OR account2_id = :account1_id
    );