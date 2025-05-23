// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: message.sql

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

const createMessage = `-- name: CreateMessage :exec
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
         )
`

type CreateMessageParams struct {
	NotifyType MessagesNotifyType
	MsgType    MessagesMsgType
	MsgContent string
	MsgExtend  json.RawMessage
	FileID     sql.NullInt64
	AccountID  sql.NullInt64
	RlyMsgID   sql.NullInt64
	RelationID int64
}

func (q *Queries) CreateMessage(ctx context.Context, arg *CreateMessageParams) error {
	_, err := q.exec(ctx, q.createMessageStmt, createMessage,
		arg.NotifyType,
		arg.MsgType,
		arg.MsgContent,
		arg.MsgExtend,
		arg.FileID,
		arg.AccountID,
		arg.RlyMsgID,
		arg.RelationID,
	)
	return err
}

const createMessageReturn = `-- name: CreateMessageReturn :one
SELECT
    id, msg_content, COALESCE(msg_extend,'{}'), file_id, create_at
FROM messages
WHERE id = LAST_INSERT_ID()
`

type CreateMessageReturnRow struct {
	ID         int64
	MsgContent string
	MsgExtend  json.RawMessage
	FileID     sql.NullInt64
	CreateAt   time.Time
}

func (q *Queries) CreateMessageReturn(ctx context.Context) (*CreateMessageReturnRow, error) {
	row := q.queryRow(ctx, q.createMessageReturnStmt, createMessageReturn)
	var i CreateMessageReturnRow
	err := row.Scan(
		&i.ID,
		&i.MsgContent,
		&i.MsgExtend,
		&i.FileID,
		&i.CreateAt,
	)
	return &i, err
}

const getLastMessageByRelation = `-- name: GetLastMessageByRelation :one
SELECT id,msg_type,msg_content,create_at
FROM messages
WHERE relation_id = ?  -- 替换为目标relation_id（如100）
  AND is_revoke = FALSE  -- 可选：排除已撤回消息
ORDER BY create_at DESC, id DESC
LIMIT 1
`

type GetLastMessageByRelationRow struct {
	ID         int64
	MsgType    MessagesMsgType
	MsgContent string
	CreateAt   time.Time
}

func (q *Queries) GetLastMessageByRelation(ctx context.Context, relationID int64) (*GetLastMessageByRelationRow, error) {
	row := q.queryRow(ctx, q.getLastMessageByRelationStmt, getLastMessageByRelation, relationID)
	var i GetLastMessageByRelationRow
	err := row.Scan(
		&i.ID,
		&i.MsgType,
		&i.MsgContent,
		&i.CreateAt,
	)
	return &i, err
}

const getMessageByID = `-- name: GetMessageByID :one
select id, notify_type, msg_type, msg_content, COALESCE(msg_extend, '{}') AS msg_extend, file_id, account_id,
       rly_msg_id, relation_id, create_at, is_revoke, is_top, is_pin, pin_time, read_ids
from messages
where id = ?
limit 1
`

type GetMessageByIDRow struct {
	ID         int64
	NotifyType MessagesNotifyType
	MsgType    MessagesMsgType
	MsgContent string
	MsgExtend  json.RawMessage
	FileID     sql.NullInt64
	AccountID  sql.NullInt64
	RlyMsgID   sql.NullInt64
	RelationID int64
	CreateAt   time.Time
	IsRevoke   bool
	IsTop      bool
	IsPin      bool
	PinTime    time.Time
	ReadIds    json.RawMessage
}

func (q *Queries) GetMessageByID(ctx context.Context, id int64) (*GetMessageByIDRow, error) {
	row := q.queryRow(ctx, q.getMessageByIDStmt, getMessageByID, id)
	var i GetMessageByIDRow
	err := row.Scan(
		&i.ID,
		&i.NotifyType,
		&i.MsgType,
		&i.MsgContent,
		&i.MsgExtend,
		&i.FileID,
		&i.AccountID,
		&i.RlyMsgID,
		&i.RelationID,
		&i.CreateAt,
		&i.IsRevoke,
		&i.IsTop,
		&i.IsPin,
		&i.PinTime,
		&i.ReadIds,
	)
	return &i, err
}

const getMsgByRelationIDAndTime = `-- name: GetMsgByRelationIDAndTime :many
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
limit ? offset ?
`

type GetMsgByRelationIDAndTimeParams struct {
	RelationID   int64
	RelationID_2 int64
	CreateAt     time.Time
	Limit        int32
	Offset       int32
}

type GetMsgByRelationIDAndTimeRow struct {
	ID         int64
	NotifyType MessagesNotifyType
	MsgType    MessagesMsgType
	MsgContent string
	MsgExtend  json.RawMessage
	FileID     sql.NullInt64
	AccountID  sql.NullInt64
	RlyMsgID   sql.NullInt64
	RelationID int64
	CreateAt   time.Time
	IsRevoke   bool
	IsTop      bool
	IsPin      bool
	PinTime    time.Time
	ReadIds    json.RawMessage
	Total      interface{}
	ReplyCount int64
}

func (q *Queries) GetMsgByRelationIDAndTime(ctx context.Context, arg *GetMsgByRelationIDAndTimeParams) ([]*GetMsgByRelationIDAndTimeRow, error) {
	rows, err := q.query(ctx, q.getMsgByRelationIDAndTimeStmt, getMsgByRelationIDAndTime,
		arg.RelationID,
		arg.RelationID_2,
		arg.CreateAt,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetMsgByRelationIDAndTimeRow{}
	for rows.Next() {
		var i GetMsgByRelationIDAndTimeRow
		if err := rows.Scan(
			&i.ID,
			&i.NotifyType,
			&i.MsgType,
			&i.MsgContent,
			&i.MsgExtend,
			&i.FileID,
			&i.AccountID,
			&i.RlyMsgID,
			&i.RelationID,
			&i.CreateAt,
			&i.IsRevoke,
			&i.IsTop,
			&i.IsPin,
			&i.PinTime,
			&i.ReadIds,
			&i.Total,
			&i.ReplyCount,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMsgsByContent = `-- name: GetMsgsByContent :many
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
LIMIT ? OFFSET ?
`

type GetMsgsByContentParams struct {
	Column1 int64
	Limit   int32
	Offset  int32
}

type GetMsgsByContentRow struct {
	ID         int64
	NotifyType MessagesNotifyType
	MsgType    MessagesMsgType
	MsgContent string
	MsgExtend  json.RawMessage
	FileID     sql.NullInt64
	AccountID  sql.NullInt64
	RelationID int64
	CreateAt   time.Time
	Total      interface{}
}

func (q *Queries) GetMsgsByContent(ctx context.Context, arg *GetMsgsByContentParams) ([]*GetMsgsByContentRow, error) {
	rows, err := q.query(ctx, q.getMsgsByContentStmt, getMsgsByContent, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetMsgsByContentRow{}
	for rows.Next() {
		var i GetMsgsByContentRow
		if err := rows.Scan(
			&i.ID,
			&i.NotifyType,
			&i.MsgType,
			&i.MsgContent,
			&i.MsgExtend,
			&i.FileID,
			&i.AccountID,
			&i.RelationID,
			&i.CreateAt,
			&i.Total,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMsgsByContentAndRelation = `-- name: GetMsgsByContentAndRelation :many
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
LIMIT ? OFFSET ?
`

type GetMsgsByContentAndRelationParams struct {
	RelationID int64
	AccountID  int64
	Limit      int32
	Offset     int32
}

type GetMsgsByContentAndRelationRow struct {
	ID         int64
	NotifyType MessagesNotifyType
	MsgType    MessagesMsgType
	MsgContent string
	MsgExtend  json.RawMessage
	FileID     sql.NullInt64
	AccountID  sql.NullInt64
	RelationID int64
	CreateAt   time.Time
	Total      interface{}
}

func (q *Queries) GetMsgsByContentAndRelation(ctx context.Context, arg *GetMsgsByContentAndRelationParams) ([]*GetMsgsByContentAndRelationRow, error) {
	rows, err := q.query(ctx, q.getMsgsByContentAndRelationStmt, getMsgsByContentAndRelation,
		arg.RelationID,
		arg.AccountID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetMsgsByContentAndRelationRow{}
	for rows.Next() {
		var i GetMsgsByContentAndRelationRow
		if err := rows.Scan(
			&i.ID,
			&i.NotifyType,
			&i.MsgType,
			&i.MsgContent,
			&i.MsgExtend,
			&i.FileID,
			&i.AccountID,
			&i.RelationID,
			&i.CreateAt,
			&i.Total,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPinMsgsByRelationID = `-- name: GetPinMsgsByRelationID :many
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
LIMIT ? OFFSET ?
`

type GetPinMsgsByRelationIDParams struct {
	RelationID   int64
	RelationID_2 int64
	Limit        int32
	Offset       int32
}

type GetPinMsgsByRelationIDRow struct {
	ID         int64
	NotifyType MessagesNotifyType
	MsgType    MessagesMsgType
	MsgContent string
	MsgExtend  json.RawMessage
	FileID     sql.NullInt64
	AccountID  sql.NullInt64
	RelationID int64
	CreateAt   time.Time
	IsRevoke   bool
	IsTop      bool
	IsPin      bool
	PinTime    time.Time
	ReadIds    json.RawMessage
	ReplyCount int64
	Total      interface{}
}

func (q *Queries) GetPinMsgsByRelationID(ctx context.Context, arg *GetPinMsgsByRelationIDParams) ([]*GetPinMsgsByRelationIDRow, error) {
	rows, err := q.query(ctx, q.getPinMsgsByRelationIDStmt, getPinMsgsByRelationID,
		arg.RelationID,
		arg.RelationID_2,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetPinMsgsByRelationIDRow{}
	for rows.Next() {
		var i GetPinMsgsByRelationIDRow
		if err := rows.Scan(
			&i.ID,
			&i.NotifyType,
			&i.MsgType,
			&i.MsgContent,
			&i.MsgExtend,
			&i.FileID,
			&i.AccountID,
			&i.RelationID,
			&i.CreateAt,
			&i.IsRevoke,
			&i.IsTop,
			&i.IsPin,
			&i.PinTime,
			&i.ReadIds,
			&i.ReplyCount,
			&i.Total,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRlyMsgsInfoByMsgID = `-- name: GetRlyMsgsInfoByMsgID :many
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
LIMIT ? OFFSET ?
`

type GetRlyMsgsInfoByMsgIDParams struct {
	RelationID   int64
	RelationID_2 int64
	Column3      int64
	Limit        int32
	Offset       int32
}

type GetRlyMsgsInfoByMsgIDRow struct {
	ID         int64
	NotifyType MessagesNotifyType
	MsgType    MessagesMsgType
	MsgContent string
	MsgExtend  json.RawMessage
	FileID     sql.NullInt64
	AccountID  sql.NullInt64
	RelationID int64
	CreateAt   time.Time
	IsRevoke   bool
	IsTop      bool
	IsPin      bool
	PinTime    time.Time
	ReadIds    json.RawMessage
	ReplyCount int64
	Total      interface{}
}

func (q *Queries) GetRlyMsgsInfoByMsgID(ctx context.Context, arg *GetRlyMsgsInfoByMsgIDParams) ([]*GetRlyMsgsInfoByMsgIDRow, error) {
	rows, err := q.query(ctx, q.getRlyMsgsInfoByMsgIDStmt, getRlyMsgsInfoByMsgID,
		arg.RelationID,
		arg.RelationID_2,
		arg.Column3,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetRlyMsgsInfoByMsgIDRow{}
	for rows.Next() {
		var i GetRlyMsgsInfoByMsgIDRow
		if err := rows.Scan(
			&i.ID,
			&i.NotifyType,
			&i.MsgType,
			&i.MsgContent,
			&i.MsgExtend,
			&i.FileID,
			&i.AccountID,
			&i.RelationID,
			&i.CreateAt,
			&i.IsRevoke,
			&i.IsTop,
			&i.IsPin,
			&i.PinTime,
			&i.ReadIds,
			&i.ReplyCount,
			&i.Total,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTopMsgByRelationID = `-- name: GetTopMsgByRelationID :one
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
LIMIT 1
`

type GetTopMsgByRelationIDParams struct {
	RelationID   int64
	RelationID_2 int64
}

type GetTopMsgByRelationIDRow struct {
	ID         int64
	NotifyType MessagesNotifyType
	MsgType    MessagesMsgType
	MsgContent string
	MsgExtend  json.RawMessage
	FileID     sql.NullInt64
	AccountID  sql.NullInt64
	RelationID int64
	CreateAt   time.Time
	IsRevoke   bool
	IsTop      bool
	IsPin      bool
	PinTime    time.Time
	ReadIds    json.RawMessage
	ReplyCount int64
	Total      interface{}
}

func (q *Queries) GetTopMsgByRelationID(ctx context.Context, arg *GetTopMsgByRelationIDParams) (*GetTopMsgByRelationIDRow, error) {
	row := q.queryRow(ctx, q.getTopMsgByRelationIDStmt, getTopMsgByRelationID, arg.RelationID, arg.RelationID_2)
	var i GetTopMsgByRelationIDRow
	err := row.Scan(
		&i.ID,
		&i.NotifyType,
		&i.MsgType,
		&i.MsgContent,
		&i.MsgExtend,
		&i.FileID,
		&i.AccountID,
		&i.RelationID,
		&i.CreateAt,
		&i.IsRevoke,
		&i.IsTop,
		&i.IsPin,
		&i.PinTime,
		&i.ReadIds,
		&i.ReplyCount,
		&i.Total,
	)
	return &i, err
}

const offerMsgsByAccountIDAndTime = `-- name: OfferMsgsByAccountIDAndTime :many
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
LIMIT ? OFFSET ?
`

type OfferMsgsByAccountIDAndTimeParams struct {
	Column1  int64
	Column2  json.RawMessage
	Column3  int64
	CreateAt time.Time
	Limit    int32
	Offset   int32
}

type OfferMsgsByAccountIDAndTimeRow struct {
	ID         int64
	NotifyType MessagesNotifyType
	MsgType    MessagesMsgType
	MsgContent string
	MsgExtend  json.RawMessage
	FileID     sql.NullInt64
	AccountID  sql.NullInt64
	RlyMsgID   sql.NullInt64
	RelationID int64
	CreateAt   time.Time
	IsRevoke   bool
	IsTop      bool
	IsPin      bool
	PinTime    time.Time
	ReadIds    json.RawMessage
	Total      interface{}
	ReplyCount int64
	HasRead    sql.NullBool
}

func (q *Queries) OfferMsgsByAccountIDAndTime(ctx context.Context, arg *OfferMsgsByAccountIDAndTimeParams) ([]*OfferMsgsByAccountIDAndTimeRow, error) {
	rows, err := q.query(ctx, q.offerMsgsByAccountIDAndTimeStmt, offerMsgsByAccountIDAndTime,
		arg.Column1,
		arg.Column2,
		arg.Column3,
		arg.CreateAt,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*OfferMsgsByAccountIDAndTimeRow{}
	for rows.Next() {
		var i OfferMsgsByAccountIDAndTimeRow
		if err := rows.Scan(
			&i.ID,
			&i.NotifyType,
			&i.MsgType,
			&i.MsgContent,
			&i.MsgExtend,
			&i.FileID,
			&i.AccountID,
			&i.RlyMsgID,
			&i.RelationID,
			&i.CreateAt,
			&i.IsRevoke,
			&i.IsTop,
			&i.IsPin,
			&i.PinTime,
			&i.ReadIds,
			&i.Total,
			&i.ReplyCount,
			&i.HasRead,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateMsgPin = `-- name: UpdateMsgPin :exec
update messages
set is_pin = ?
where id = ?
`

type UpdateMsgPinParams struct {
	IsPin bool
	ID    int64
}

func (q *Queries) UpdateMsgPin(ctx context.Context, arg *UpdateMsgPinParams) error {
	_, err := q.exec(ctx, q.updateMsgPinStmt, updateMsgPin, arg.IsPin, arg.ID)
	return err
}

const updateMsgReads = `-- name: UpdateMsgReads :exec

UPDATE messages
SET read_ids = CASE
                   WHEN JSON_CONTAINS(read_ids, CAST(@accountID AS JSON)) = 0
                       THEN JSON_ARRAY_APPEND(read_ids, '$', @accountID)
                   ELSE read_ids
    END
WHERE
    relation_id = ?
  AND JSON_CONTAINS(@msgIDs, CAST(id AS JSON))
  AND JSON_CONTAINS(read_ids, CAST(@accountID AS JSON)) = 0
`

// -- name: UpdateMsgReads :exec
// UPDATE messages AS m
// SET read_ids =
//
//	CASE
//	    WHEN JSON_CONTAINS(read_ids, CAST(@accountID AS JSON)) = 0
//	        THEN JSON_ARRAY_APPEND(read_ids, '$', @accountID)
//	    ELSE read_ids
//	    END
//
// WHERE relation_id = ?
//
//	AND JSON_CONTAINS(@target_ids, CAST(m.id AS JSON));
//
// 先执行更新操作
func (q *Queries) UpdateMsgReads(ctx context.Context, relationID int64) error {
	_, err := q.exec(ctx, q.updateMsgReadsStmt, updateMsgReads, relationID)
	return err
}

const updateMsgReadsReturn = `-- name: UpdateMsgReadsReturn :many

SELECT id, CAST(@accountID AS UNSIGNED) AS account_id,relation_id
FROM messages
WHERE
    relation_id = ?
  AND JSON_CONTAINS(@msgIDs, CAST(id AS JSON))
  AND JSON_SEARCH(read_ids, 'one', CAST(@accountID AS CHAR)) IS NOT NULL
`

type UpdateMsgReadsReturnRow struct {
	ID         int64
	AccountID  int64
	RelationID int64
}

// 再查询受影响的行
func (q *Queries) UpdateMsgReadsReturn(ctx context.Context, relationID int64) ([]*UpdateMsgReadsReturnRow, error) {
	rows, err := q.query(ctx, q.updateMsgReadsReturnStmt, updateMsgReadsReturn, relationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*UpdateMsgReadsReturnRow{}
	for rows.Next() {
		var i UpdateMsgReadsReturnRow
		if err := rows.Scan(&i.ID, &i.AccountID, &i.RelationID); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateMsgRevoke = `-- name: UpdateMsgRevoke :exec
UPDATE messages
SET is_revoke = ?
WHERE id = ?
`

type UpdateMsgRevokeParams struct {
	IsRevoke bool
	ID       int64
}

func (q *Queries) UpdateMsgRevoke(ctx context.Context, arg *UpdateMsgRevokeParams) error {
	_, err := q.exec(ctx, q.updateMsgRevokeStmt, updateMsgRevoke, arg.IsRevoke, arg.ID)
	return err
}

const updateMsgTop = `-- name: UpdateMsgTop :exec
UPDATE messages
SET is_top = ?
WHERE id = ?
`

type UpdateMsgTopParams struct {
	IsTop bool
	ID    int64
}

func (q *Queries) UpdateMsgTop(ctx context.Context, arg *UpdateMsgTopParams) error {
	_, err := q.exec(ctx, q.updateMsgTopStmt, updateMsgTop, arg.IsTop, arg.ID)
	return err
}
