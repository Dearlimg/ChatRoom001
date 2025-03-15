package tx

import (
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/dao/redis/operate"
	"context"
)

func (store *SqlStore) AddSettingWithTX(ctx context.Context, rdb *operate.RDB, accountID, relationID int64, isLeader bool) error {
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
