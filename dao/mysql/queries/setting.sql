-- name: CreateSetting :exec
insert into settings (account_id, relation_id, nick_name,is_leader,is_self)
values (?,?,'',?,?);

-- name: DeleteSetting :exec
delete
from settings
where account_id = ?
and relation_id = ?;

-- name: DeleteSettingsByAccountID :exec
delete
from settings
where account_id =?;

-- name: GetRelationIDByAccountIDFromSettings :many
select relation_id
FROM settings
where account_id = ?;

-- name: ExistGroupLeaderByAccountIDWithLock :one
select exists(select 1
              from settings
              where account_id =?
              and is_leader=true);

-- name: UpdateSettingNickName :exec
update settings
set nick_name = ?
where account_id =?
and relation_id =?;

-- name: UpdateSettingPin :exec
update settings
set is_pin = ?
where account_id=?
and  relation_id= ?;

-- name: UpdateSettingDisturb :exec
update settings
set is_not_disturb = ?
where account_id = ?
and relation_id =?;

-- name: UpdateSettingShow :exec
update settings
set  is_show = ?
where account_id =?
and relation_id =?;

-- name: GetSettingByID :one
select *
from settings
where account_id = ?
and relation_id =?;


-- name: GetFriendPinSettingsOrderByPinTime :many
select s.*,
       a.id as account_id,
       a.name as account_name,
       a.avatar as account_avatar
from (select settings.relation_id, settings.nick_name, settings.pin_time,settings.is_pin,settings.is_show,settings.is_not_disturb
      from settings,
           relations
      where settings.account_id = ?
        and settings.is_pin = true
        and settings.relation_id = relations.id
        and relations.relation_type = 'friend') as s,
     accounts a
where a.id = (select account_id from settings where relation_id = s.relation_id and (settings.account_id != ? or is_self = true))
order by s.pin_time;

-- name: GetGroupPinSettingsOrderByPinTime :many
select s.relation_id,
       s.nick_name,
       s.pin_time,
       s.is_show,
       s.is_not_disturb,
       s.is_pin,
        r.id,
        r.group_name,
        r.group_description,
        r.group_avatar,
        r.group_name,
        r.group_description,
        r.group_avatar
from (select settings.relation_id,settings.nick_name,settings.pin_time,settings.is_pin,settings.is_show,settings.is_not_disturb
      from settings,
           relations
      where settings.account_id = ?
      and settings.relation_id=relations.id
      and settings.is_pin=true
      and relation_type='group') as s,
    relations r
where r.id=(select relation_id
            from settings
            where relation_id = s.relation_id
            and settings.account_id = ?)
order by s.pin_time;

-- name: GetGroupShowSettingsOrderByShowTime :many
select s.*,
       r.id,r.group_avatar,r.group_name,r.group_description,r.created_at
from (select relation_id,
             nick_name,
             is_not_disturb,
             is_pin,
             pin_time,
             is_show,
             last_show,
             is_self
      from settings,
           relations
      where settings.account_id = ?
        and settings.relation_id = relations.id
        and settings.is_show = true
        and relations.relation_type = 'group') as s,
     relations r
where r.id = (select relation_id from settings where relation_id = s.relation_id and settings.account_id = ? limit 1)
order by s.last_show desc;

-- name: GetFriendShowSettingsOrderByShowTime :many
select s.*,
       a.id as account_id,
       a.name as account_name,
       a.avatar as account_avatar
from (select relation_id,
             nick_name,
             is_not_disturb,
             is_pin,
             pin_time,
             is_show,
             last_show,
             is_self
      from settings,
           relations
      where settings.account_id = ?
        and settings.is_show = true
        and settings.relation_id = relations.id
        and relations.relation_type = 'friend') as s,
     accounts a
where a.id = (select account_id from settings where relation_id = s.relation_id and (settings.account_id != ? or is_self = true))
order by s.last_show desc;

-- name: GetFriendSettingsOrderByName :many
select s.*,
       a.id as account_id,
       a.name as account_name,
       a.avatar as account_avatar
from (select relation_id,
             nick_name,
             is_not_disturb,
             is_pin,
             pin_time,
             is_show,
             last_show,
             is_self
      from settings,
           relations
      where settings.account_id = ?
        and settings.relation_id = relations.id
        and relations.relation_type = 'friend') as s,
     accounts a
where a.id = (select account_id from settings where relation_id = s.relation_id and (settings.account_id != ? or s.is_self = true))
order by a.name;

-- name: ExistsFriendSetting :one
SELECT EXISTS (
    SELECT 1
    FROM settings s, relations r
    WHERE
        r.relation_type = 'friend'
      AND (
        (r.account1_id = ? AND r.account2_id = ?)
            OR
        (r.account1_id = ? AND r.account2_id = ?)
        )
      AND s.account_id = ?
);


-- name: GetFriendSettingsByName :many
SELECT
    s.*,
    a.id AS account_id,
    a.name AS account_name,
    a.avatar AS account_avatar,
    COUNT(*) OVER () AS total
FROM (
         SELECT
             s.relation_id,
             s.nick_name,
             s.is_not_disturb,
             s.is_pin,
             s.pin_time,
             s.is_show,
             s.last_show,
             s.is_self
         FROM settings s
                  INNER JOIN relations r ON s.relation_id = r.id
         WHERE
             s.account_id = ?
           AND r.relation_type = 'friend'
     ) AS s
         CROSS JOIN accounts a
WHERE
    a.id = (
        SELECT account_id
        FROM settings
        WHERE
            relation_id = s.relation_id
          AND (settings.account_id != ? OR s.is_self = 1)
    )
  AND (
    a.name LIKE CONCAT('%', ?, '%')
        OR s.nick_name LIKE CONCAT('%', ?, '%')
    )
ORDER BY a.name
LIMIT ? OFFSET ?;

-- name: GetGroupSettingsByName :many
select s.*,
       r.id as realtion_id,
       r.group_name AS group_name,
       r.group_avatar AS group_avatar,
       r.group_description AS description,
       count(*) over () as total
from (select relation_id,
    nick_name,
    is_not_disturb,
    is_pin,
    pin_time,
    is_show,
    last_show,
    is_self
    from settings,
    relations
    where settings.account_id = ?
    and settings.relation_id = relations.id
    and relations.relation_type = 'group') as s,
    relations r
where r.id = (select s.relation_id from settings where s.relation_id=s.relation_id and (settings.account_id=?))
    and ((r.group_name like ('%' || ? || '%')))
order by (r.group_name)
limit ? offset ?;


-- name: TransferIsLeaderTrue :exec
UPDATE settings
SET is_leader = 1
WHERE relation_id = ? AND account_id = ?;

-- name: TransferIsLeaderFalse :exec
UPDATE settings
SET is_leader = 0
WHERE relation_id = ? AND account_id = ?;

-- name: DeleteGroup :exec
DELETE FROM settings
WHERE relation_id = ?;

-- name: ExistsSetting :one
SELECT EXISTS (
    SELECT 1
    FROM settings
    WHERE account_id = ? AND relation_id = ?
);

-- name: ExistsIsLeader :one
SELECT EXISTS (
    SELECT 1
    FROM settings
    WHERE relation_id = ? AND account_id = ? AND is_leader = 1
);

-- name: GetGroupMembers :many
SELECT account_id
FROM settings
WHERE relation_id = ?;

-- name: GetAccountIDsByRelationID :many
SELECT DISTINCT account_id
FROM settings
WHERE relation_id = ?;

-- name: GetGroupList :many
select s.*,
       r.id as relation_id,
       r.group_name as group_name,
       r.group_description as group_discription,
       r.group_avatar as group_avatar,
       count(*) over () as total
from (select relation_id,
             nick_name,
             is_not_disturb,
             is_pin,
             pin_time,
             is_show,
             last_show,
             is_self
      from settings,
           relations
      where settings.account_id = ?
        and settings.relation_id = relations.id
        and relations.relation_type = 'group') as s,
    relations r
where r.id = (select s.relation_id from settings where relation_id = s.relation_id and (settings.account_id = ?))
order by s.last_show;

-- name: CreateManySetting :exec
INSERT INTO settings (account_id, relation_id, nick_name)
VALUES (?, ?, ?);

-- name: GetGroupMembersByID :many
SELECT
    a.id,
    a.name,
    a.avatar,
    s.nick_name,
    s.is_leader
FROM accounts a
         LEFT JOIN settings s ON a.id = s.account_id
WHERE s.relation_id = ?
LIMIT ? OFFSET ?;