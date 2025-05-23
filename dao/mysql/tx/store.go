package tx

import (
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/dao/redis/operate"
	"ChatRoom001/pkg/tool"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type TXer interface {
	//CreateAccountWithTx 创建账号并建立和自己的关系
	CreateAccountWithTx(ctx context.Context, rdb *operate.RDB, maxAccountNum int64, arg *db.CreateAccountParams) error
	//DeleteAccountWithTx 删除账号并删除与之相关的好友关系
	DeleteAccountWithTx(ctx context.Context, rdb *operate.RDB, accountID int64) error
	//CreateApplicationTx 判断是否存在申请，不存在则创建申请
	CreateApplicationTx(ctx context.Context, params *db.CreateApplicationParams) error
	//AcceptApplicationTx account2 接受 account1 的申请并建立好友关系和双方的关系设置，同时发送消息通知并添加到 redis
	AcceptApplicationTx(ctx context.Context, rdb *operate.RDB, account1, account2 *db.GetAccountByIDRow) (*db.Message, error)
	//DeleteRelationWithTx 从数据库中删除关系并删除 redis 中的关系
	DeleteRelationWithTx(ctx context.Context, rdb *operate.RDB, relationID int64) error
	//AddSettingWithTx 向数据库和 redis 中同时添加群成员
	AddSettingWithTx(ctx context.Context, rdb *operate.RDB, accountID, relationID int64, isLeader bool) error
	// UploadGroupAvatarWithTx 创建群组头像文件
	UploadGroupAvatarWithTx(ctx context.Context, arg db.CreateFileParams) error
	// TransferGroupWithTx 转让群
	TransferGroupWithTx(ctx context.Context, accountID, relationID, toAccountID int64) error
	// DeleteSettingWithTx 从数据库和 redis 中删除群员
	DeleteSettingWithTx(ctx context.Context, rdb *operate.RDB, accountID, relationID int64) error
	// RevokeMsgWithTx 撤回消息，如果消息 pin 或者置顶，则全部取消
	RevokeMsgWithTx(ctx context.Context, msgID int64, isPin, isTop bool) error
}

// sqlStore 处理数据库操作
type SqlStore struct {
	*db.Queries // 嵌入 *db.Queries，可以直接访问 db.Queries 中定义的方法和字段，不需要间接访问
	Pool        *sql.DB
}

func (store *SqlStore) AddSettingWithTx(ctx context.Context, rdb *operate.RDB, accountID, relationID int64, isLeader bool) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		err := queries.CreateSetting(ctx, &db.CreateSettingParams{
			AccountID:  accountID,
			RelationID: relationID,
			IsLeader:   isLeader,
			IsSelf:     false,
		})
		if err != nil {
			return err
		}
		return rdb.AddRelationAccount(ctx, relationID, accountID)
	})
}

func (store *SqlStore) UploadGroupAvatarWithTx(ctx context.Context, arg db.CreateFileParams) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		_, err = queries.GetGroupAvatar(ctx, arg.RelationID)
		if err != nil {
			// 如果没有设置过群头像
			//if err.Error() == "no rows in result set" {
			if errors.Is(sql.ErrNoRows, err) {
				err = queries.CreateFile(ctx, &db.CreateFileParams{
					FileName:   arg.FileName,
					FileType:   "img",
					FileSize:   arg.FileSize,
					Key:        arg.Key,
					Url:        arg.Url,
					RelationID: arg.RelationID,
					AccountID:  sql.NullInt64{},
				})
			} else {
				return err
			}
		} else {
			// 在 file 表中覆盖之前的群头像
			err = queries.UpdateGroupAvatar(ctx, &db.UpdateGroupAvatarParams{
				Url:        arg.Url,
				RelationID: arg.RelationID,
			})
		}
		data, err := queries.GetGroupRelationByID(ctx, arg.RelationID.Int64)
		if err != nil {
			return err
		}
		// 更新 relation 表中的头像数据
		return queries.UpdateGroupRelation(ctx, &db.UpdateGroupRelationParams{
			GroupName: data.GroupName,

			GroupDescription: data.GroupDescription,
			GroupAvatar: sql.NullString{
				String: arg.Url,
				Valid:  true,
			},
			ID: arg.RelationID.Int64,
		})
	})
}

func (store *SqlStore) TransferGroupWithTx(ctx context.Context, accountID, relationID, toAccountID int64) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		err = tool.DoThat(err, func() error {
			// 将原群主的 isLeader 转换为 false
			return queries.TransferIsLeaderFalse(ctx, &db.TransferIsLeaderFalseParams{
				RelationID: relationID,
				AccountID:  accountID,
			})
		})
		err = tool.DoThat(err, func() error {
			// 将新群主的 isLeader 转换为 true
			return queries.TransferIsLeaderTrue(ctx, &db.TransferIsLeaderTrueParams{
				RelationID: relationID,
				AccountID:  toAccountID,
			})
		})
		return err
	})
}

func (store *SqlStore) DeleteSettingWithTx(ctx context.Context, rdb *operate.RDB, accountID, relationID int64) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		err := queries.DeleteSetting(ctx, &db.DeleteSettingParams{
			AccountID:  accountID,
			RelationID: relationID,
		})
		if err != nil {
			return err
		}
		return rdb.DeleteRelationAccount(ctx, relationID, accountID)
	})
}

func (store *SqlStore) RevokeMsgWithTx(ctx context.Context, msgID int64, isPin, isTop bool) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		err = tool.DoThat(err, func() error {
			return queries.UpdateMsgRevoke(ctx, &db.UpdateMsgRevokeParams{
				ID:       msgID,
				IsRevoke: true,
			})
		})
		if isPin {
			err = tool.DoThat(err, func() error {
				return queries.UpdateMsgPin(ctx, &db.UpdateMsgPinParams{
					ID:    msgID,
					IsPin: false,
				})
			})
		}
		if isTop {
			err = tool.DoThat(err, func() error {
				return queries.UpdateMsgTop(ctx, &db.UpdateMsgTopParams{
					ID:    msgID,
					IsTop: false,
				})
			})
		}
		return err
	})
}

func (store *SqlStore) execTx(ctx context.Context, fn func(queries *db.Queries) error) error {
	// 开启数据库事务
	tx, err := store.Pool.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted, // 设置事务隔离级别为已提交读
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// 创建一个包含事务的查询对象（你可以在这里传递事务对象，查询会在事务中执行）
	q := store.WithTx(tx)
	// 执行传入的回调函数
	if err := fn(q); err != nil {
		// 如果回调执行失败，回滚事务
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// 如果回调执行成功，提交事务
	return tx.Commit()
}
