package tx

import (
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/dao/redis/operate"
	"context"
)

func (store *SqlStore) DeleteRelationWithTx(ctx context.Context, rdb *operate.RDB, relationID int64) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		err := queries.DeleteRelation(ctx, relationID)
		if err != nil {
			return err
		}
		return rdb.DeleteRelations(ctx, relationID)
	})
}
