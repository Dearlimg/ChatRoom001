package tx

import (
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/dao/redis/operate"
	"context"
)

type TXer interface {
	// CreateAccountWithTx 创建账号并建立和自己的关系
	CreateAccountWithTx(ctx context.Context, rdb *operate.RDB, maxAccountNum int32, arg *db.CreateAccountParams) error
	// DeleteAccountWithTx 删除账号并删除与之相关的好友关系
	DeleteAccountWithTx(ctx context.Context, rdb *operate.RDB, accountID int64) error
	// CreateApplicationTx 判断是否存在申请，不存在则创建申请
	//CreateApplicationTx(ctx context.Context, params *db.CreateApplicationParams) error
	// AcceptApplicationTx account2 接受 account1 的申请并建立好友关系和双方的关系设置，同时发送消息通知并添加到 redis
	//AcceptApplicationTx(ctx context.Context, rdb *operate.RDB, account1, account2 *db.GetAccountByIDRow) (*db.Message, error)
	// DeleteRelationWithTx 从数据库中删除关系并删除 redis 中的关系
	DeleteRelationWithTx(ctx context.Context, rdb *operate.RDB, relationID int64) error
	// AddSettingWithTx 向数据库和 redis 中同时添加群成员
	AddSettingWithTx(ctx context.Context, rdb *operate.RDB, accountID, relationID int64, isLeader bool) error
	// UploadGroupAvatarWithTx 创建群组头像文件
	//UploadGroupAvatarWithTx(ctx context.Context, arg db.CreateFileParams) error
	// TransferGroupWithTx 转让群
	TransferGroupWithTx(ctx context.Context, accountID, relationID, toAccountID int64) error
	// DeleteSettingWithTx 从数据库和 redis 中删除群员
	DeleteSettingWithTx(ctx context.Context, rdb *operate.RDB, accountID, relationID int64) error
	// RevokeMsgWithTx 撤回消息，如果消息 pin 或者置顶，则全部取消
	RevokeMsgWithTx(ctx context.Context, msgID int64, isPin, isTop bool) error
}
