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
		//var relationID int64

		ID := sql.NullInt64{
			Int64: arg.ID,
			Valid: true,
		}

		err = tool.DoThat(err, func() error {
			err = queries.CreateFriendRelation(ctx, &db.CreateFriendRelationParams{
				Account1ID: ID,
				Account2ID: ID,
			})
			return err
		})
		err = tool.DoThat(err, func() error {
			return queries.CreateSt
		})
	})

}
