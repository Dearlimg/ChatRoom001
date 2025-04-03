package logic

import (
	"ChatRoom001/dao"
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/reply"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type message struct{}

func GetMsgInfoByID(ctx context.Context, msgID int64) (*db.GetMessageByIDRow, errcode.Err) {
	result, err := dao.Database.DB.GetMessageByID(ctx, msgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.MsgNotExists
		}
		global.Logger.Error(err.Error())
		return nil, errcode.ErrServer
	}
	return result, nil
}

func (message) CreateFileMsg(ctx *gin.Context, params model.CreateFileMsg) (*reply.ParamCreateFileMsg, errcode.Err) {
	ok, myErr := ExistsSetting(ctx, params.AccountID, params.RelationID)
	if myErr != nil {
		return nil, myErr
	}
	if !ok {
		return nil, errcodes.AuthPermissionsInsufficient
	}
	//fileInfo, myErr := Logics.File.PublishFile()
	return &reply.ParamCreateFileMsg{}, nil
}

func (message) GetMsgsByRelationIDAndTime(ctx *gin.Context, params model.GetMsgsByRelationIDAndTime) (*reply.ParamGetMsgsRelationIDAndTime, errcode.Err) {
	ok, myErr := ExistsSetting(ctx, params.AccountID, params.RelationID)
	fmt.Println("GetMsgsByRelationIDAndTime : ", params.AccountID, params.RelationID)
	if myErr != nil {
		return nil, myErr
	}
	if !ok {
		return nil, errcodes.AuthPermissionsInsufficient
	}
	data, err := dao.Database.DB.GetMsgByRelationIDAndTime(ctx, &db.GetMsgByRelationIDAndTimeParams{
		RelationID:   params.RelationID,
		RelationID_2: params.RelationID,
		CreateAt:     params.LastTime,
		Limit:        params.Limit,
		Offset:       params.Offset,
	})
	fmt.Println("GetMsgsByRelationIDAndTime1 : ", data)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if len(data) == 0 {
		return &reply.ParamGetMsgsRelationIDAndTime{
			List: []*reply.ParamMsgInfoWithRly{},
		}, nil
	}
	result := make([]*reply.ParamMsgInfoWithRly, len(data))
	for _, v := range data {
		var content string
		var extend *model.MsgExtend
		if !v.IsRevoke { // 该消息没有被撤回
			content = v.MsgContent
			extend, err = model.JsonToExtend(v.MsgExtend)
			if err != nil {
				global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
				continue
			}
		}
		var readIDs []int64
		if params.AccountID == v.AccountID.Int64 {
			if err := json.Unmarshal(v.ReadIds, &readIDs); err != nil {
				return nil, errcode.ErrServer
			}
		}
		var rlyMsg *reply.ParamRlyMsg
		if v.RlyMsgID.Valid { // 该 ID 有意义
			rlyMsgInfo, myErr := GetMsgInfoByID(ctx, v.RlyMsgID.Int64)
			if myErr != nil {
				continue
			}
			var rlyContent string
			var rlyExtend *model.MsgExtend
			if !rlyMsgInfo.IsRevoke { // 回复消息没有撤回
				rlyContent = rlyMsgInfo.MsgContent
				rlyExtend, err = model.JsonToExtend(rlyMsgInfo.MsgExtend)
				if err != nil {
					global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
					continue
				}
			}
			rlyMsg = &reply.ParamRlyMsg{
				MsgID:      v.RlyMsgID.Int64,
				MsgType:    string(rlyMsgInfo.MsgType),
				MsgContent: rlyContent,
				MsgExtend:  rlyExtend,
				IsRevoked:  rlyMsgInfo.IsRevoke,
			}
		}
		result = append(result, &reply.ParamMsgInfoWithRly{
			ParamMsgInfo: reply.ParamMsgInfo{
				ID:         v.ID,
				NotifyType: string(v.NotifyType),
				MsgType:    string(v.MsgType),
				MsgContent: content,
				MsgExtend:  extend,
				FileID:     v.FileID.Int64,
				AccountID:  v.AccountID.Int64,
				RelationID: v.RelationID,
				CreateAt:   v.CreateAt,
				IsRevoke:   v.IsRevoke,
				IsTop:      v.IsTop,
				IsPin:      v.IsPin,
				PinTime:    v.PinTime,
				ReadIds:    readIDs,
				ReplyCount: v.ReplyCount,
			},
			RlyMsg: rlyMsg,
		})
	}
	return &reply.ParamGetMsgsRelationIDAndTime{List: result, Total: data[0].Total.(int64)}, nil
}

func (message) OfferMsgsByAccountIDAndTime(ctx *gin.Context, params model.OfferMsgsByAccountIDAndTime) (*reply.ParamOfferMsgsByAccountIDAndTime, errcode.Err) {
	data, err := dao.Database.DB.OfferMsgsByAccountIDAndTime(ctx, &db.OfferMsgsByAccountIDAndTimeParams{
		CreateAt: params.LastTime,
		Limit:    params.Limit,
		Offset:   params.Offset,
		Column3:  params.AccountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if len(data) == 0 {
		return &reply.ParamOfferMsgsByAccountIDAndTime{List: []*reply.ParamMsgInfoWithRlyAndHasRead{}}, nil
	}
	result := make([]*reply.ParamMsgInfoWithRlyAndHasRead, 0, len(data))
	for _, v := range data {
		var content string
		var extend *model.MsgExtend
		if !v.IsRevoke {
			content = v.MsgContent
			extend, err = model.JsonToExtend(v.MsgExtend)
			if err != nil {
				global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
				continue
			}
		}
		//var readIDs []int64
		//if params.AccountID == v.AccountID.Int64 {
		//	readIDs = v.ReadIds
		//}
		var readIDs []int64
		if params.AccountID == v.AccountID.Int64 {
			if err := json.Unmarshal(v.ReadIds, &readIDs); err != nil {
				return nil, errcode.ErrServer
			}
		}
		var rlyMsg *reply.ParamRlyMsg
		if v.RlyMsgID.Valid {
			rlyMsgInfo, myErr := GetMsgInfoByID(ctx, v.RlyMsgID.Int64)
			if myErr != nil {
				continue
			}
			var rlyContent string
			var rlyExtend *model.MsgExtend
			if !rlyMsgInfo.IsRevoke {
				rlyContent = rlyMsgInfo.MsgContent
				rlyExtend, err = model.JsonToExtend(rlyMsgInfo.MsgExtend)
				if err != nil {
					global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
					continue
				}
			}
			rlyMsg = &reply.ParamRlyMsg{
				MsgID:      rlyMsgInfo.ID,
				MsgType:    string(rlyMsgInfo.MsgType),
				MsgContent: rlyContent,
				MsgExtend:  rlyExtend,
				IsRevoked:  rlyMsgInfo.IsRevoke,
			}
		}
		result = append(result, &reply.ParamMsgInfoWithRlyAndHasRead{
			ParamMsgInfoWithRly: reply.ParamMsgInfoWithRly{
				ParamMsgInfo: reply.ParamMsgInfo{
					ID:         v.ID,
					NotifyType: string(v.NotifyType),
					MsgType:    string(v.MsgType),
					MsgContent: content,
					MsgExtend:  extend,
					FileID:     v.FileID.Int64,
					AccountID:  v.AccountID.Int64,
					RelationID: v.RelationID,
					CreateAt:   v.CreateAt,
					IsRevoke:   v.IsRevoke,
					IsTop:      v.IsTop,
					IsPin:      v.IsPin,
					PinTime:    v.PinTime,
					ReadIds:    readIDs,
					ReplyCount: v.ReplyCount,
				},
				RlyMsg: rlyMsg,
			},
			HasRead: false, //v.HasRead,
		})
	}
	return &reply.ParamOfferMsgsByAccountIDAndTime{List: result, Total: data[0].Total.(int64)}, nil
}

func (message) GetPinMsgsByRelationID(ctx *gin.Context, accountID, relationID int64, limit, offset int32) (*reply.ParamGetPinMsgsByRelationID, errcode.Err) {
	ok, err := ExistsSetting(ctx, accountID, relationID)
	if err != nil {
		return &reply.ParamGetPinMsgsByRelationID{Total: 0}, err
	}
	if !ok {
		return &reply.ParamGetPinMsgsByRelationID{Total: 0}, errcodes.AuthPermissionsInsufficient
	}
	data, myerr := dao.Database.DB.GetPinMsgsByRelationID(ctx, &db.GetPinMsgsByRelationIDParams{
		RelationID: relationID,
		Limit:      limit,
		Offset:     offset,
	})
	if myerr != nil {
		global.Logger.Error(myerr.Error(), middlewares.ErrLogMsg(ctx)...)
		return &reply.ParamGetPinMsgsByRelationID{Total: 0}, errcode.ErrServer
	}
	if len(data) == 0 {
		return &reply.ParamGetPinMsgsByRelationID{List: []*reply.ParamMsgInfo{}}, nil
	}
	result := make([]*reply.ParamMsgInfo, 0, len(data))
	for _, v := range data {
		var content string
		var extend *model.MsgExtend
		if !v.IsRevoke {
			content = v.MsgContent
			extend, myerr = model.JsonToExtend(v.MsgExtend)
			if myerr != nil {
				global.Logger.Error(myerr.Error(), middlewares.ErrLogMsg(ctx)...)
				return &reply.ParamGetPinMsgsByRelationID{Total: 0}, errcode.ErrServer
			}
		}
		var readIDs []int64
		if accountID == v.AccountID.Int64 {
			if err := json.Unmarshal(v.ReadIds, &readIDs); err != nil {
				return nil, errcode.ErrServer
			}
		}
		//var readIDs []int64
		//if accountID == v.AccountID.Int64 {
		//	readIDs = v.ReadIds
		//}
		result = append(result, &reply.ParamMsgInfo{
			ID:         v.ID,
			NotifyType: string(v.NotifyType),
			MsgType:    string(v.MsgType),
			MsgContent: content,
			MsgExtend:  extend,
			FileID:     v.FileID.Int64,
			AccountID:  v.AccountID.Int64,
			RelationID: v.RelationID,
			CreateAt:   v.CreateAt,
			IsRevoke:   v.IsRevoke,
			IsTop:      v.IsTop,
			IsPin:      v.IsPin,
			PinTime:    v.PinTime,
			ReadIds:    readIDs,
			ReplyCount: v.ReplyCount,
		})
	}
	return &reply.ParamGetPinMsgsByRelationID{
		List:  result,
		Total: data[0].Total.(int64),
	}, nil
}
