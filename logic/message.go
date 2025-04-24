package logic

import (
	"ChatRoom001/dao"
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/chat/server"
	format2 "ChatRoom001/model/format"
	"ChatRoom001/model/reply"
	"ChatRoom001/model/request"
	"ChatRoom001/task"
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
	fileInfo, myErr := Logics.File.PublishFile(ctx, model.PublishFile{
		File:       params.File,
		RelationID: params.RelationID,
		AccountID:  params.AccountID,
	})

	//fileType, myErr := gtype.GetFileType(params.File)

	if myErr != nil {
		return nil, myErr
	}
	var isRly bool  // 是否是回复别人的消息
	var rlyID int64 // 回复 ID 为 rlyID 的消息
	var rlyMsg *reply.ParamRlyMsg
	if params.RlyMsgID > 0 { // 如果是回复别人的消息
		rltInfo, myErr := GetMsgInfoByID(ctx, params.RlyMsgID)
		if myErr != nil {
			return nil, myErr
		}
		if rltInfo.IsRevoke {
			return nil, errcodes.RlyMsgHasRevoked
		}
		isRly = true
		rlyID = params.RlyMsgID
		rlyMsgExtend, err := model.JsonToExtend(rltInfo.MsgExtend)
		if err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return nil, errcode.ErrServer
		}
		rlyMsg = &reply.ParamRlyMsg{
			MsgID:   rltInfo.ID,
			MsgType: string(rltInfo.MsgType),
			//MsgType:    fileType,
			MsgContent: rltInfo.MsgContent,
			MsgExtend:  rlyMsgExtend,
			IsRevoked:  rltInfo.IsRevoke,
		}
	}
	extend, _ := model.ExtendToJson(nil)
	err := dao.Database.DB.CreateMessage(ctx, &db.CreateMessageParams{
		NotifyType: db.MessagesNotifyTypeCommon,
		MsgType:    db.MessagesMsgType(model.MsgTypeFile),
		//MsgType:    db.MessagesMsgType(fileType),
		MsgContent: fileInfo.Url,
		MsgExtend:  extend,
		FileID:     sql.NullInt64{Int64: fileInfo.ID, Valid: true},
		AccountID:  sql.NullInt64{Int64: params.AccountID, Valid: true},
		RlyMsgID:   sql.NullInt64{Int64: rlyID, Valid: isRly},
		RelationID: params.RelationID,
	})

	result, err := dao.Database.DB.CreateMessageReturn(ctx)

	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}

	global.Worker.SendTask(task.PublishMsg(reply.ParamMsgInfoWithRly{
		ParamMsgInfo: reply.ParamMsgInfo{
			ID:         result.ID,
			NotifyType: string(db.MessagesNotifyTypeCommon),
			MsgType:    string(model.MsgTypeFile),
			MsgContent: result.MsgContent,
			MsgExtend:  nil,
			AccountID:  params.AccountID,
			RelationID: params.RelationID,
			CreateAt:   result.CreateAt,
		},
		RlyMsg: rlyMsg,
	}))
	return &reply.ParamCreateFileMsg{
		ID:         result.ID,
		MsgContent: result.MsgContent,
		FileID:     result.FileID.Int64,
		CreateAt:   result.CreateAt,
	}, nil

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
	result := make([]*reply.ParamMsgInfoWithRly, 0, len(data))
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
		RelationID:   relationID,
		RelationID_2: relationID,
		Limit:        limit,
		Offset:       offset,
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

func (message) GetRlyMsgsInfoByMsgID(ctx *gin.Context, accountID, relationID, msgID int64, limit, offset int32) (*reply.ParamGetRlyMsgsInfoByMsgID, errcode.Err) {
	ok, err := ExistsSetting(ctx, accountID, relationID)
	if err != nil {
		return &reply.ParamGetRlyMsgsInfoByMsgID{}, err
	}
	if !ok {
		return &reply.ParamGetRlyMsgsInfoByMsgID{}, errcodes.AuthPermissionsInsufficient
	}
	data, myErr := dao.Database.DB.GetRlyMsgsInfoByMsgID(ctx, &db.GetRlyMsgsInfoByMsgIDParams{
		RelationID:   relationID,
		RelationID_2: relationID,
		Limit:        limit,
		Offset:       offset,
		Column3:      msgID,
	})
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return &reply.ParamGetRlyMsgsInfoByMsgID{Total: 0}, errcode.ErrServer
	}
	if len(data) == 0 {
		return &reply.ParamGetRlyMsgsInfoByMsgID{List: []*reply.ParamMsgInfo{}}, nil
	}
	result := make([]*reply.ParamMsgInfo, 0, len(data))
	for _, v := range data {
		var content string
		var extend *model.MsgExtend
		if !v.IsRevoke {
			content = v.MsgContent
			extend, myErr = model.JsonToExtend(v.MsgExtend)
			if myErr != nil {
				global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
				return &reply.ParamGetRlyMsgsInfoByMsgID{Total: 0}, errcode.ErrServer
			}
		}
		var readIDs []int64
		if accountID == v.AccountID.Int64 {
			if err := json.Unmarshal(v.ReadIds, &readIDs); err != nil {
				return nil, errcode.ErrServer
			}
		}
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
	return &reply.ParamGetRlyMsgsInfoByMsgID{
		List:  result,
		Total: data[0].Total.(int64),
	}, nil
}

func (message) GetMsgsByContent(ctx *gin.Context, accountID, relationID int64, content string, limit, offset int32) (*reply.ParamGetMsgsByContent, errcode.Err) {
	return &reply.ParamGetMsgsByContent{}, errcode.ErrServer
}

func (message) GetTopMsgByRelationID(ctx *gin.Context, accountID, relationID int64) (*reply.ParamGetTopMsgByRelationID, errcode.Err) {
	ok, err := ExistsSetting(ctx, accountID, relationID)
	if err != nil {
		return &reply.ParamGetTopMsgByRelationID{}, err
	}
	if !ok {
		return &reply.ParamGetTopMsgByRelationID{}, errcodes.AuthPermissionsInsufficient
	}
	data, myerr := dao.Database.DB.GetTopMsgByRelationID(ctx, &db.GetTopMsgByRelationIDParams{
		RelationID:   relationID,
		RelationID_2: relationID,
	})
	if myerr != nil {
		if errors.Is(myerr, sql.ErrNoRows) {
			return nil, nil
		}
		global.Logger.Error(myerr.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	var content string
	var extend *model.MsgExtend
	if !data.IsRevoke {
		content = data.MsgContent
		extend, myerr = model.JsonToExtend(data.MsgExtend)
		if myerr != nil {
			global.Logger.Error(myerr.Error(), middlewares.ErrLogMsg(ctx)...)
			return nil, errcode.ErrServer
		}
	}
	var readIDs []int64
	if accountID == relationID {
		if err := json.Unmarshal(data.ReadIds, &readIDs); err != nil {
			return nil, errcode.ErrServer
		}
	}
	fmt.Println("GetTopMsgByRelationID  ", data.AccountID.Int64)
	return &reply.ParamGetTopMsgByRelationID{MsgInfo: reply.ParamMsgInfo{
		ID:         data.ID,
		NotifyType: string(data.NotifyType),
		MsgType:    string(data.MsgType),
		MsgContent: content,
		MsgExtend:  extend,
		FileID:     data.FileID.Int64,
		AccountID:  data.AccountID.Int64,
		RelationID: data.RelationID,
		CreateAt:   data.CreateAt,
		IsRevoke:   data.IsRevoke,
		IsTop:      data.IsTop,
		IsPin:      data.IsPin,
		PinTime:    data.PinTime,
		ReadIds:    readIDs,
		ReplyCount: data.ReplyCount,
	}}, nil
}

func (message) UpdateMsgPin(ctx *gin.Context, accountID int64, param *request.ParamUpdateMsgPin) errcode.Err {
	ok, err := ExistsSetting(ctx, accountID, param.RelationID)
	if err != nil {
		return err
	}

	if !ok {
		return errcodes.AuthPermissionsInsufficient
	}
	msgInfo, err := GetMsgInfoByID(ctx, param.ID)
	if err != nil {
		return err
	}
	if msgInfo.IsPin == param.IsPin {
		return nil
	}

	myerr := dao.Database.DB.UpdateMsgPin(ctx, &db.UpdateMsgPinParams{
		ID:    param.ID,
		IsPin: param.IsPin,
	})

	if myerr != nil {
		global.Logger.Error(myerr.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}

	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.UpdateMsgState(accessToken, param.RelationID, param.ID, server.MsgPin, param.IsPin))

	return nil
}

func (message) UpdateMsgTop(ctx *gin.Context, accountID int64, param *request.ParamUpdateMsgTop) errcode.Err {
	ok, err := ExistsSetting(ctx, accountID, param.RelationID)
	if err != nil {
		return err
	}
	if !ok {
		return errcodes.AuthPermissionsInsufficient
	}
	msgInfo, err := GetMsgInfoByID(ctx, param.ID)
	if err != nil {
		return err
	}
	if msgInfo.IsTop == param.IsTop {
		return nil
	}

	msg, _ := dao.Database.DB.GetTopMsgByRelationID(ctx, &db.GetTopMsgByRelationIDParams{
		RelationID:   param.RelationID,
		RelationID_2: param.RelationID,
	})

	myerr := dao.Database.DB.UpdateMsgTop(ctx, &db.UpdateMsgTopParams{
		ID:    msg.ID,
		IsTop: !msg.IsTop,
	})

	myerr = dao.Database.DB.UpdateMsgTop(ctx, &db.UpdateMsgTopParams{
		ID:    param.ID,
		IsTop: param.IsTop,
	})
	if myerr != nil {
		global.Logger.Error(myerr.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)

	global.Worker.SendTask(task.UpdateMsgState(accessToken, param.RelationID, param.ID, server.MsgTop, param.IsTop))
	//fmt.Println("UpdateMsgTop ", param.RelationID, param.ID, server.MsgTop, param.IsTop)
	f := func() error {

		name, err := dao.Database.DB.GetAccountNameByID(ctx, accountID)
		if err != nil {
			global.Logger.Error(myerr.Error(), middlewares.ErrLogMsg(ctx)...)
			return errcode.ErrServer
		}

		var tempMsg string
		if param.IsTop == true {
			tempMsg = fmt.Sprintf(format2.TopMessage, name)
		} else {
			tempMsg = fmt.Sprintf(format2.UnTopMessage, name)
		}

		arg := &db.CreateMessageParams{
			NotifyType: db.MessagesNotifyTypeSystem,
			MsgType:    db.MessagesMsgType(model.MsgTypeText),
			MsgContent: tempMsg,
			//MsgContent: fmt.Sprintf(format2.TopMessage, accountID),
			RelationID: msgInfo.RelationID,
		}
		err = dao.Database.DB.CreateMessage(ctx, arg)
		if err != nil {
			return err
		}
		msgRly, err := dao.Database.DB.CreateMessageReturn(ctx)
		global.Worker.SendTask(task.PublishMsg(reply.ParamMsgInfoWithRly{
			ParamMsgInfo: reply.ParamMsgInfo{
				ID:         msgRly.ID,
				NotifyType: string(arg.NotifyType),
				MsgType:    string(arg.MsgType),
				MsgContent: arg.MsgContent,
				RelationID: arg.RelationID,
				CreateAt:   msgRly.CreateAt,
				AccountID:  -1,
			},
			RlyMsg: nil,
		}))
		return nil
	}
	if err := f(); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		reTry("UpdateMsgTop", f)
	}
	return nil
}

func (message) RevokeMsg(ctx *gin.Context, accountID int64, msgID int64) errcode.Err {
	msgInfo, err := GetMsgInfoByID(ctx, msgID)
	if err != nil {
		return err
	}
	if msgInfo.AccountID.Int64 != accountID {
		return errcodes.AuthPermissionsInsufficient
	}
	if msgInfo.IsRevoke {
		return errcodes.MsgAlreadyRevoke
	}
	myErr := dao.Database.DB.RevokeMsgWithTx(ctx, msgID, msgInfo.IsPin, msgInfo.IsTop)
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.UpdateMsgState(accessToken, msgInfo.RelationID, msgID, server.MsgRevoke, true))
	if msgInfo.IsTop {
		// 推送 top 通知
		global.Worker.SendTask(task.UpdateMsgState(accessToken, msgInfo.RelationID, msgID, server.MsgTop, false))
		// 创建并推送 top 消息
		f := func() error {

			arg := &db.CreateMessageParams{
				NotifyType: db.MessagesNotifyTypeSystem,
				MsgType:    db.MessagesMsgType(model.MsgTypeText),
				MsgContent: fmt.Sprintf(format2.TopMessage, accountID),
				RelationID: msgInfo.RelationID,
			}
			err := dao.Database.DB.CreateMessage(ctx, arg)
			msgRly, err := dao.Database.DB.CreateMessageReturn(ctx)
			if err != nil {
				return err
			}
			global.Worker.SendTask(task.PublishMsg(reply.ParamMsgInfoWithRly{
				ParamMsgInfo: reply.ParamMsgInfo{
					ID:         msgRly.ID,
					NotifyType: string(arg.NotifyType),
					MsgType:    string(arg.MsgType),
					MsgContent: arg.MsgContent,
					RelationID: arg.RelationID,
					CreateAt:   msgRly.CreateAt,
				},
				RlyMsg: nil,
			}))
			return nil
		}
		if err := f(); err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			reTry("RevokeMsg", f)
		}
	}
	return nil
}
