package tx

import (
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/dao/redis/operate"
	"ChatRoom001/errcodes"
	"ChatRoom001/pkg/tool"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

func (store *SqlStore) CreateApplicationTx(ctx context.Context, param *db.CreateApplicationParams) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		ok, err := queries.ExistsApplicationByIDWithLock(ctx, &db.ExistsApplicationByIDWithLockParams{
			Account1ID: param.Account1ID,
			Account2ID: param.Account2ID,
		})
		//fmt.Println("application tx 16", ok, err, param.Account1ID, param.Account2ID, param.ApplyMsg)
		if err != nil {
			return err
		}
		if ok {
			return errcodes.ApplicationExists
		}
		return queries.CreateApplication(ctx, param)
	})
}

func (store *SqlStore) AcceptApplicationTx(ctx context.Context, rdb *operate.RDB, account1, account2 *db.GetAccountByIDRow) (*db.Message, error) {
	var result *db.Message
	fmt.Println("AcceptApplicationTx is hrr ")
	err := store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		err = tool.DoThat(err, func() error {
			return queries.UpdateApplication(ctx, &db.UpdateApplicationParams{
				Status:     db.ApplicationsStatusValue1,
				Account1ID: account1.ID,
				Account2ID: account2.ID,
			})
		})
		id1, id2 := account1.ID, account2.ID
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		var relationID int64

		var tempid1 sql.NullInt64
		tempid1.Valid = true
		tempid1.Int64 = id1
		var tempid2 sql.NullInt64
		tempid2.Valid = true
		tempid2.Int64 = id2

		err = tool.DoThat(err, func() error {
			err = queries.CreateFriendRelation(ctx, &db.CreateFriendRelationParams{
				Account1ID: tempid1,
				Account2ID: tempid2,
			})
			return err
		})

		relationID, err = queries.CreateRelationReturn(ctx)
		// 建立双方关系
		err = tool.DoThat(err, func() error {
			return queries.CreateSetting(ctx, &db.CreateSettingParams{
				AccountID:  account1.ID,
				RelationID: relationID,
				IsLeader:   false,
				IsSelf:     false,
			})
		})
		err = tool.DoThat(err, func() error {
			return queries.CreateSetting(ctx, &db.CreateSettingParams{
				AccountID:  account2.ID,
				RelationID: relationID,
				IsLeader:   false,
				IsSelf:     false,
			})
		})
		// 新建一个系统通知消息作为好友的第一条消息
		err = tool.DoThat(err, func() error {
			arg := &db.CreateMessageParams{
				NotifyType: db.MessagesNotifyTypeCommon,
				MsgType:    db.MessagesMsgTypeText,
				MsgContent: "我们已经成为好友啦，现在可以开始聊天啦！",
				MsgExtend:  json.RawMessage{},
				AccountID:  sql.NullInt64{Int64: account2.ID, Valid: true},
				RelationID: relationID,
			}
			err := queries.CreateMessage(ctx, arg)
			msgInfo, err := queries.CreateMessageReturn(ctx)
			result = &db.Message{
				ID:         msgInfo.ID,
				NotifyType: arg.NotifyType,
				MsgType:    arg.MsgType,
				MsgContent: arg.MsgContent,
				RelationID: relationID,
				CreateAt:   msgInfo.CreateAt,
			}
			return err
		})
		err = tool.DoThat(err, func() error {
			return rdb.AddRelationAccount(ctx, relationID, account1.ID, account2.ID)
		})
		return err
	})
	return result, err
}
