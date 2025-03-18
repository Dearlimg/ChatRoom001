package tx

import (
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/errcodes"
	"context"
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

//func (store *SqlStore)AcceptApplicationTx(ctx context.Context,rdb *operate.RDB,account1,account2 *db.GetAccountByIDRow)(*db.Message, error)
