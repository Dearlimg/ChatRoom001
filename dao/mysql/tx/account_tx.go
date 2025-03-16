package tx

import (
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/dao/redis/operate"
	"ChatRoom001/pkg/tool"
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

var (
	ErrAccountOverNum     = errors.New("账户数量超过限制")
	ErrAccountNameExist   = errors.New("账户名已经存在")
	ErrAccountGroupLeader = errors.New("账号是群主")
)

// CreateAccountWithTx 检查数量、账户名之后创建账户并建立和自己的关系
func (store *SqlStore) CreateAccountWithTx(ctx context.Context, rdb *operate.RDB, maxAccountNum int64, arg *db.CreateAccountParams) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		var accountNum int64
		// 检查数量
		err = tool.DoThat(err, func() error {
			accountNum, err = queries.CountAccountByUserID(ctx, arg.UserID)
			return err
		})
		if accountNum >= maxAccountNum {
			return ErrAccountOverNum
		}
		// 检查账户名
		var exists bool
		err = tool.DoThat(err, func() error {
			exists, err = queries.ExistsAccountByNameAndUserID(ctx, &db.ExistsAccountByNameAndUserIDParams{
				UserID: arg.UserID,
				Name:   arg.Name,
			})
			return err
		})
		if exists {
			return ErrAccountNameExist
		}
		// 创建账户
		err = tool.DoThat(err, func() error {
			return queries.CreateAccount(ctx, arg)
		})
		// 建立关系(自己与自己的好友关系)
		ID := sql.NullInt64{
			Int64: arg.UserID,
			Valid: true,
		}

		err = tool.DoThat(err, func() error {
			err = queries.CreateFriendRelation(ctx, &db.CreateFriendRelationParams{
				Account1ID: ID,
				Account2ID: ID,
			})
			return err
		})
		ID1 := sql.NullInt64{Int64: arg.UserID, Valid: true} // arg.ID 是账户的 ID（来自 accounts 表）
		err = queries.CreateFriendRelation(ctx, &db.CreateFriendRelationParams{
			Account1ID: ID,
			Account2ID: ID,
		})

		param1 := db.GetRelationIDByInfoParams{
			Account1ID: ID1,
			Account2ID: ID1,
		}
		rID, err := queries.GetRelationIDByInfo(ctx, &param1)

		err = tool.DoThat(err, func() error {
			return queries.CreateSetting(ctx, &db.CreateSettingParams{
				AccountID:  arg.ID,
				RelationID: rID,
				IsSelf:     true,
			})
		})
		err = tool.DoThat(err, func() error { return rdb.AddRelationAccount(ctx, rID) })
		return err
	})
}

func (store *SqlStore) DeleteAccountWithTx(ctx context.Context, rdb *operate.RDB, accountID int64) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		// 判断该账户是否是群主
		var isLeader bool
		err = tool.DoThat(err, func() error {
			isLeader, err = queries.ExistGroupLeaderByAccountIDWithLock(ctx, accountID)
			return err
		})
		if isLeader {
			return ErrAccountGroupLeader
		}
		// 删除好友

		var friendRelationIDs []int64
		err = tool.DoThat(err, func() error {
			err = queries.DeleteSettingsByAccountID(ctx, accountID)
			return err
		})

		//groupRelationID,err:=queries.GetRelationIDByAccountIDFromSettings(ctx,accountID)

		// 删除群
		var groupRelationIDs []int64
		err = tool.DoThat(err, func() error {
			err = queries.DeleteSettingsByAccountID(ctx, accountID)
			return err
		})
		// 删除账户
		err = tool.DoThat(err, func() error {
			err = queries.DeleteAccount(ctx, accountID)
			return err
		})
		// 从 redis 中删除对应的关系
		// 从 redis 中删除该账户的好友关系
		err = tool.DoThat(err, func() error {
			return rdb.DeleteRelations(ctx, friendRelationIDs...)
		})
		// 在 redis 中删除该账户所在的群聊中的该账户
		err = tool.DoThat(err, func() error {
			return rdb.DeleteAccountFromRelations(ctx, accountID, groupRelationIDs...)
		})
		return err
	})
}
