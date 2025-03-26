package logic

import (
	"ChatRoom001/dao"
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/middlewares"
	"ChatRoom001/task"
	"database/sql"
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type application struct {
}

func (application) CreateApplication(ctx *gin.Context, accountID1, accountID2 int64, msg string) errcode.Err {
	if accountID1 == accountID2 {
		return errcodes.ApplicationNotValid
	}
	id1, id2 := sortID(accountID1, accountID2)

	var ID1 sql.NullInt64
	ID1.Int64 = id1

	var ID2 sql.NullInt64
	ID2.Int64 = id2

	exist, err := dao.Database.DB.ExistsFriendRelation(ctx, &db.ExistsFriendRelationParams{
		Account2ID: ID2,
		Account1ID: ID1,
	})

	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}

	if exist {
		return errcodes.RelationExists
	}
	//fmt.Println("logic application 45", accountID1, accountID2)
	err = dao.Database.DB.CreateApplicationTx(ctx, &db.CreateApplicationParams{
		Account1ID: accountID1,
		Account2ID: accountID2,
		ApplyMsg:   msg,
	})

	switch {
	case errors.Is(err, errcodes.ApplicationExists):
		return errcodes.ApplicationExists
	case errors.Is(err, nil):
		// 提示对方有新地申请消息
		global.Worker.SendTask(task.Application(accountID2))
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
}

func (application) DeleteApplication(ctx *gin.Context, accountID1, accountID2 int64) errcode.Err {
	apply, err := getApplication(ctx, accountID1, accountID2)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if apply.Account1ID != accountID1 {
		return errcodes.AuthPermissionsInsufficient
	}
	if err := dao.Database.DB.DeleteApplication(ctx, &db.DeleteApplicationParams{
		Account1ID: accountID1,
		Account2ID: accountID2,
	}); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	return nil
}

func getApplication(ctx *gin.Context, accountID1, accountID2 int64) (*db.Application, errcode.Err) {
	apply, err := dao.Database.DB.GetApplicationByID(ctx, &db.GetApplicationByIDParams{
		Account1ID: accountID1,
		Account2ID: accountID2,
	})
	switch {
	case errors.Is(err, nil):
		return apply, nil
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errcodes.ApplicationNotExists
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
}

func (application) RefuseApplication(ctx *gin.Context, accountID1, accountID2 int64, refuseMsg string) errcode.Err {
	apply, err := getApplication(ctx, accountID1, accountID2)
	if err != nil {
		return err
	}
	if apply.Status == db.ApplicationsStatusValue2 {
		return errcodes.ApplicationRepeatOpt
	}
	if err1 := dao.Database.DB.UpdateApplication(ctx, &db.UpdateApplicationParams{
		RefuseMsg:  refuseMsg,
		Status:     db.ApplicationsStatusValue2,
		Account1ID: accountID1,
		Account2ID: accountID2,
	}); err1 != nil {
		global.Logger.Error(err1.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	return nil
}

func (application) AcceptApplication(ctx *gin.Context, accountID1, accountID2 int64) errcode.Err {
	apply, err := getApplication(ctx, accountID1, accountID2)
	if err != nil {
		return err
	}
	if apply.Status == db.ApplicationsStatusValue1 {
		return errcodes.ApplicationRepeatOpt
	}
	accountInfo1, myerr := GetAccountInfoByID(ctx, accountID1, accountID2)
	if myerr != nil {
		return myerr
	}
	accountInfo2, myerr := GetAccountInfoByID(ctx, accountID1, accountID2)
	if myerr != nil {
		return myerr
	}
	fmt.Println(accountInfo1, accountInfo2)
	//msgInfo, err := dao.Database.DB.AcceptApplicationTx(ctx, dao.Database.Redis, accountInfo1, accountInfo2)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	return nil
}

func (application) ListApplications(ctx *gin.Context) {}
