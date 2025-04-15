package logic

import (
	"ChatRoom001/dao"
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/middlewares"
	"ChatRoom001/task"
	"database/sql"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type group struct{}

func (group) CreateGroup(ctx *gin.Context, accountID int64, name, description string) (relationID int64, err errcode.Err) {
	var tempname sql.NullString
	tempname.String = name
	tempname.Valid = true
	var tempdescr sql.NullString
	tempdescr.String = description
	tempdescr.Valid = true
	var tempavater sql.NullString
	tempavater.String = global.PublicSetting.Rules.DefaultAvatarURL
	tempavater.Valid = true

	myErr := dao.Database.DB.CreateGroupRelation(ctx, &db.CreateGroupRelationParams{
		GroupName:        tempname,
		GroupDescription: tempdescr,
		GroupAvatar:      tempavater,
	})

	relationID, newErr := dao.Database.DB.CreateGroupRelationReturn(ctx, &db.CreateGroupRelationReturnParams{
		GroupName:        tempname,
		GroupDescription: tempdescr,
	})

	if newErr != nil {
		global.Logger.Error(myErr.Error())
		return 0, errcode.ErrServer
	}

	if myErr != nil {
		global.Logger.Error(myErr.Error())
		return 0, errcode.ErrServer
	}
	myErr = dao.Database.DB.AddSettingWithTx(ctx, dao.Database.Redis, accountID, relationID, true)
	if myErr != nil {
		global.Logger.Error(myErr.Error())
		return 0, errcode.ErrServer
	}
	return relationID, nil
}

func (group) TransferGroup(ctx *gin.Context, accountID, relationID, toAccountID int64) errcode.Err {
	ok, err := dao.Database.DB.ExistsIsLeader(ctx, &db.ExistsIsLeaderParams{
		RelationID: relationID,
		AccountID:  accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if !ok {
		return errcodes.NotLeader
	}
	ok, err = dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if !ok {
		return errcodes.NotGroupMember
	}
	err = dao.Database.DB.TransferGroupWithTx(ctx, accountID, relationID, toAccountID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	// 推送群主更换的通知
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.TransferGroup(accessToken, accountID, relationID))
	return nil
}
