-- name: CreateMessage :exec
INSERT INTO messages (
    notify_type,
    msg_type,
    msg_content,
    msg_extend,
    file_id,
    account_id,
    rly_msg_id,
    relation_id,
    read_ids
) VALUES (
             ?, ?, ?, ?, ?, ?, ?, ?, JSON_ARRAY()
         );

-- name: CreateMessageReturn :one
SELECT
    id, msg_content, COALESCE(msg_extend,'{}'), file_id, create_at
FROM messages
WHERE id = LAST_INSERT_ID();

-- name: GetMessageByID :one
select id, notify_type, msg_type, msg_content, COALESCE(msg_extend, '{}') AS msg_extend, file_id, account_id,
       rly_msg_id, relation_id, create_at, is_revoke, is_top, is_pin, pin_time, read_ids
from messages
where id = ?
limit 1;


# -- name: UpdateMsgReads :exec
# UPDATE messages AS m
# SET read_ids =
#         CASE
#             WHEN JSON_CONTAINS(read_ids, CAST(@accountID AS JSON)) = 0
#                 THEN JSON_ARRAY_APPEND(read_ids, '$', @accountID)
#             ELSE read_ids
#             END
# WHERE relation_id = ?
#   AND JSON_CONTAINS(@target_ids, CAST(m.id AS JSON));
-- 先执行更新操作

-- name: UpdateMsgReads :exec
UPDATE messages
SET read_ids = CASE
                   WHEN JSON_CONTAINS(read_ids, CAST(@accountID AS JSON)) = 0
                       THEN JSON_ARRAY_APPEND(read_ids, '$', @accountID)
                   ELSE read_ids
    END
WHERE
    relation_id = ?
  AND JSON_CONTAINS(@msgIDs, CAST(id AS JSON))
  AND JSON_CONTAINS(read_ids, CAST(@accountID AS JSON)) = 0;

-- 再查询受影响的行

-- name: UpdateMsgReadsReturn :many
SELECT id, CAST(@accountID AS UNSIGNED) AS account_id,relation_id
FROM messages
WHERE
    relation_id = ?
  AND JSON_CONTAINS(@msgIDs, CAST(id AS JSON))
  AND JSON_SEARCH(read_ids, 'one', CAST(@accountID AS CHAR)) IS NOT NULL;

-- name: GetMsgByRelationIDAndTime :many
select m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       coalesce(m1.msg_extend,'{}'),
       m1.file_id,
       m1.account_id,
       m1.rly_msg_id,
       m1.relation_id,
       m1.create_at,
       m1.is_revoke,
       m1.is_top,
       m1.is_pin,
       m1.pin_time,
       m1.read_ids,
       count(*) over () as total,
       (select count(id) from messages where rly_msg_id = m1.id and messages.relation_id = ?) as reply_count
from messages m1
where m1.relation_id = ?
  and m1.create_at < ?
order by m1.create_at
limit ? offset ?;

-- name: OfferMsgsByAccountIDAndTime :many
SELECT m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       coalesce(m1.msg_extend,'{}'),
       m1.file_id,
       m1.account_id,
       m1.rly_msg_id,
       m1.relation_id,
       m1.create_at,
       m1.is_revoke,
       m1.is_top,
       m1.is_pin,
       m1.pin_time,
       m1.read_ids,
       COUNT(*) OVER () AS total,
       (SELECT COUNT(id) FROM messages WHERE rly_msg_id = m1.id AND relation_id = m1.relation_id) AS reply_count,
       (m1.account_id = CAST(? AS UNSIGNED) OR JSON_CONTAINS(m1.read_ids, CAST(? AS JSON), '$')) AS has_read
FROM messages m1
         JOIN settings s ON m1.relation_id = s.relation_id AND s.account_id = CAST(? AS UNSIGNED)
WHERE m1.create_at > ?
LIMIT ? OFFSET ?;

-- name: UpdateMsgPin :exec
update messages
set is_pin = ?
where id = ?;

-- name: UpdateMsgTop :exec
UPDATE messages
SET is_top = ?
WHERE id = ?;

-- name: UpdateMsgRevoke :exec
UPDATE messages
SET is_revoke = ?
WHERE id = ?;

-- name: GetTopMsgByRelationID :one
SELECT m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       coalesce(m1.msg_extend,'{}'),
       m1.file_id,
       m1.account_id,
       m1.relation_id,
       m1.create_at,
       m1.is_revoke,
       m1.is_top,
       m1.is_pin,
       m1.pin_time,
       m1.read_ids,
       (SELECT COUNT(id) FROM messages WHERE rly_msg_id = m1.id AND messages.relation_id = ?) AS reply_count,
       COUNT(*) OVER () AS total
FROM messages m1
WHERE m1.relation_id = ? AND m1.is_top = TRUE
ORDER BY pin_time DESC
LIMIT 1;

-- name: GetPinMsgsByRelationID :many
SELECT m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       coalesce(m1.msg_extend,'{}'),
       m1.file_id,
       m1.account_id,
       m1.relation_id,
       m1.create_at,
       m1.is_revoke,
       m1.is_top,
       m1.is_pin,
       m1.pin_time,
       m1.read_ids,
       (SELECT COUNT(id) FROM messages WHERE rly_msg_id = m1.id AND messages.relation_id = ?) AS reply_count,
       COUNT(*) OVER () AS total
FROM messages m1
WHERE m1.relation_id = ? AND m1.is_pin = TRUE
ORDER BY m1.pin_time DESC
LIMIT ? OFFSET ?;

-- name: GetRlyMsgsInfoByMsgID :many
SELECT m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       coalesce(m1.msg_extend,'{}'),
       m1.file_id,
       m1.account_id,
       m1.relation_id,
       m1.create_at,
       m1.is_revoke,
       m1.is_top,
       m1.is_pin,
       m1.pin_time,
       m1.read_ids,
       (SELECT COUNT(id) FROM messages WHERE rly_msg_id = m1.id AND messages.relation_id = ?) AS reply_count,
       COUNT(*) OVER () AS total
FROM messages m1
WHERE m1.relation_id = ? AND m1.rly_msg_id = CAST(? AS UNSIGNED)
ORDER BY m1.create_at
LIMIT ? OFFSET ?;

-- name: GetMsgsByContentAndRelation :many
SELECT m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       coalesce(m1.msg_extend,'{}'),
       m1.file_id,
       m1.account_id,
       m1.relation_id,
       m1.create_at,
       COUNT(*) OVER () AS total
FROM messages m1
         JOIN settings s ON m1.relation_id = ? AND m1.relation_id = s.relation_id AND s.account_id = ?
WHERE (NOT m1.is_revoke)
  AND MATCH(m1.msg_content_tsy) AGAINST(? IN NATURAL LANGUAGE MODE)
ORDER BY m1.create_at DESC
LIMIT ? OFFSET ?;


-- name: GetMsgsByContent :many
SELECT m1.id,
       m1.notify_type,
       m1.msg_type,
       m1.msg_content,
       coalesce(m1.msg_extend,'{}'),
       m1.file_id,
       m1.account_id,
       m1.relation_id,
       m1.create_at,
       COUNT(*) OVER () AS total
FROM messages m1
         JOIN settings s ON m1.relation_id = s.relation_id AND s.account_id = CAST(? AS UNSIGNED)
WHERE (NOT is_revoke)
  AND MATCH(m1.msg_content_tsy) AGAINST(? IN NATURAL LANGUAGE MODE)
ORDER BY m1.create_at DESC
LIMIT ? OFFSET ?;

-- name: GetLastMessageByRelation :one
SELECT id,msg_type,msg_content,create_at
FROM messages
WHERE relation_id = ?  -- 替换为目标relation_id（如100）
  AND is_revoke = FALSE  -- 可选：排除已撤回消息
ORDER BY create_at DESC, id DESC
LIMIT 1;