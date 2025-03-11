-- name: CreateSetting :exec
insert into settings (account_id, relation_id, nick_name,is_leader,is_self)
values (?,?,'',?,?);

-- name: DeleteSetting :exec
delete
from settings
where account_id = ?
and relation_id = ?;

-- name: DeleteSettingsByAccountID :many
delete
from settings
where account_id =?;

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
from (select settings.relation_id, settings.nick_name, settings.pin_time
      from settings,
           relations
      where settings.account_id = ?
        and settings.is_pin = true
        and settings.relation_id = relations.id
        and relations.relation_type = 'friend') as s,
     accounts a
where a.id = (select account_id from settings where relation_id = s.relation_id and (account_id != ? or is_self = true))
order by s.pin_time;

-- name: GetGroupPinSettingsOrderByPinTime :many
select s.relation_id,
       s.nick_name,
       s.pin_time,
        r.id,
        r.group_name,
        r.group_description,
        r.group_avatar
from (select settings.relation_id,settings.nick_name,settings.pin_time
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
            and account_id = ?)
order by s.pin_time;

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
where a.id = (select account_id from settings where relation_id = s.relation_id and (account_id != ? or is_self = true))
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
where a.id = (select account_id from settings where relation_id = s.relation_id and (account_id != ? or s.is_self = true))
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
