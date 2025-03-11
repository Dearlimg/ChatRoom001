-- name: CreateGroupRelation :exec
INSERT INTO relations (relation_type, group_name, group_description, group_avatar)
value ('group',?,?,?);

-- name: CreateFriendRelation :exec
insert into relations(relation_type,account1_id,account2_id)
value ('friend',?,?);

-- name: DeleteFriendRelation :exec
insert into relations(relation_type,account1_id,account2_id)
value ('friend',?,?);

-- name: DeleteRelation :exec
delete
from relations
where id=?;

-- name: DeleteFriendRelationByAccountID :exec
delete
from relations
where relation_type ='friend'
and (account1_id=?);

-- name: UpdateGroupRelation :exec
UPDATE relations
SET group_name = ?, group_description = ?, group_avatar = ?
WHERE relation_type = 'group'
AND id = ?;

-- name: GetGroupRelationByID :one
select id,relation_type,relations.group_name,relations.group_description,relations.group_avatar
from relations
where relation_type = 'group' and id = ?;

-- name: ExistsFriendRelation :one
select exists(select 1
              from relations
              where relation_type = 'friend'
              and account1_id =?
              and account2_id=?);

-- name: GetFriendRelationByID :one
select (relations.account2_id,relations.account1_id,relations.created_at)
from relations
where relation_type='friend'
  and id =?;

-- name: GetAllGroupRelation :many
select id
from relations
where relation_type = 'group'
and account1_id is null
and account2_id is null;

-- name: GetAllRelationOnRelation :many
select *
from relations;

-- name: GetAllRelationIDs :many
select id
from relations;

-- name: GetRelationIDByAccountID :one
select id
from relations
where account2_id=?
and account1_id=?;

