package logic

import (
	"ChatRoom001/dao"
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/reply"
	"ChatRoom001/model/request"
	"ChatRoom001/task"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
	"time"
)

type notify struct{}

func (notify) CreateNotify(ctx *gin.Context, accountID int64, params *request.ParamCreateNotify) (*reply.ParamGroupNotify, errcode.Err) {
	ok, err := dao.Database.DB.ExistsIsLeader(ctx, &db.ExistsIsLeaderParams{
		AccountID:  accountID,
		RelationID: params.RelationID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if !ok {
		return nil, errcodes.NotLeader
	}
	extend, _ := model.ExtendToJson(params.MsgExtend)
	readIDsJSON, err := json.Marshal([]int64{accountID})
	myErr := dao.Database.DB.CreateGroupNotify(ctx, &db.CreateGroupNotifyParams{
		RelationID: sql.NullInt64{Int64: params.RelationID, Valid: true},
		MsgContent: params.MsgContent,
		MsgExpand:  extend,
		AccountID:  sql.NullInt64{Int64: accountID, Valid: true},
		CreateAt:   time.Now(),
		ReadIds:    readIDsJSON,
	})
	if myErr != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	data, _ := dao.Database.DB.CreateGroupNotifyReturn(ctx)
	msgExtend, err := model.ExtendToJson(params.MsgExtend)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}

	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.CreateNotify(accessToken, accountID, params.RelationID, data.MsgContent, params.MsgExtend))
	return &reply.ParamGroupNotify{
		ID:         data.ID,
		RelationID: data.RelationID.Int64,
		MsgContent: data.MsgContent,
		MsgExtend:  msgExtend,
		AccountID:  data.AccountID.Int64,
		CreateAt:   data.CreateAt,
		ReadIDs:    data.ReadIds,
	}, nil
}

func (notify) UpdateNotify(ctx *gin.Context, accountID int64, params *request.ParamUpdateNotify) errcode.Err {
	ok, err := dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		AccountID:  accountID,
		RelationID: params.RelationID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if !ok {
		return errcodes.NotLeader
	}
	extend, _ := model.ExtendToJson(params.MsgExtend)
	readIDsJSON, err := json.Marshal([]int64{accountID})
	err = dao.Database.DB.UpdateGroupNotify(ctx, &db.UpdateGroupNotifyParams{
		RelationID: sql.NullInt64{Int64: params.RelationID, Valid: true},
		MsgContent: params.MsgContent,
		MsgExpand:  extend,
		AccountID:  sql.NullInt64{Int64: accountID, Valid: true},
		CreateAt:   time.Now(),
		ReadIds:    readIDsJSON,
		ID:         params.ID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.UpdateNotify(accessToken, accountID, params.RelationID, params.MsgContent, params.MsgExtend))
	return nil
}

func (notify) GetNotifyByID(ctx *gin.Context, accountID int64, params *request.ParamGetNotifyByID) (*reply.ParamGetNotifyByID, errcode.Err) {
	ok, err := dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		AccountID:  accountID,
		RelationID: params.RelationID,
	})
	fmt.Println("GetNotifyByID content2:")
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if !ok {
		return nil, errcode.ErrServer
	}

	var tempRelationID sql.NullInt64
	tempRelationID.Valid = true
	tempRelationID.Int64 = params.RelationID

	data, err := dao.Database.DB.GetGroupNotifyByID(ctx, tempRelationID)
	fmt.Println("GetNotifyByID content3:", data, err)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	result := make([]reply.ParamGroupNotify, 0, len(data))
	for _, v := range data {
		fmt.Println("GetNotifyByID4   as dawd ", v)
		if err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return nil, errcode.ErrServer
		}
		result = append(result, reply.ParamGroupNotify{
			ID:         v.ID,
			RelationID: v.RelationID.Int64,
			MsgContent: v.MsgContent,
			MsgExtend:  v.MsgExpand,
			AccountID:  v.AccountID.Int64,
			CreateAt:   v.CreateAt,
			ReadIDs:    v.ReadIds,
		})

	}
	return &reply.ParamGetNotifyByID{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

func (notify) DeleteNotify(ctx *gin.Context, accountID int64, params *request.ParamDeleteNotify) errcode.Err {
	ok, err := dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		AccountID:  accountID,
		RelationID: params.RelationID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if !ok {
		return errcode.ErrServer
	}
	err = dao.Database.DB.DeleteGroupNotify(ctx, params.ID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	return nil
}
