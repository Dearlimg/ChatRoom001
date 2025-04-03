package chat

import (
	"ChatRoom001/dao"
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/logic"
	"ChatRoom001/model"
	"ChatRoom001/model/chat/client"
	"ChatRoom001/model/reply"
	"ChatRoom001/task"
	"context"
	"database/sql"
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
)

type message struct{}

func (message) SendMsg(ctx context.Context, param *model.HandleSendMsg) (*client.HandleSendMsgRly, errcode.Err) {
	//fmt.Println("SendMsg0")
	ok, myErr := logic.ExistsSetting(ctx, param.AccountID, param.RelationID)
	if myErr != nil {
		//fmt.Println("SendMsg2", myErr)
		return nil, myErr
	}
	if !ok {
		//fmt.Println("SendMsg3", param)
		return nil, errcodes.AuthPermissionsInsufficient
	}
	var rlyMsgID int64
	var rlyMsg *reply.ParamRlyMsg
	//fmt.Println("SendMsg4")
	if param.RlyMsgID > 0 {
		rlyInfo, myErr := logic.GetMsgInfoByID(ctx, param.RlyMsgID)
		if myErr != nil {
			fmt.Println("SendMsg5", myErr)
			return nil, myErr
		}
		if rlyInfo.RelationID != param.RelationID {
			//fmt.Println("SendMsg6")
			return nil, errcodes.RlyMsgNotOneRelation
		}
		if rlyInfo.IsRevoke {
			return nil, errcodes.RlyMsgHasRevoked
		}
		rlyMsgID = param.RlyMsgID
		rlyMsgExtend, err := model.JsonToExtend(rlyInfo.MsgExtend)
		if err != nil {
			global.Logger.Error(err.Error())
			return nil, errcode.ErrServer
		}
		//fmt.Println("SendMsg7", err)
		rlyMsg = &reply.ParamRlyMsg{
			MsgID:      rlyInfo.ID,
			MsgType:    string(rlyInfo.MsgType),
			MsgContent: rlyInfo.MsgContent,
			MsgExtend:  rlyMsgExtend,
			IsRevoked:  rlyInfo.IsRevoke,
		}
	}
	//fmt.Println("SendMsg4.8")
	msgExtend, err := model.ExtendToJson(param.MsgExtend)
	if err != nil {
		global.Logger.Error(err.Error())
		return nil, errcode.ErrServer
	}
	//fmt.Println("SendMsg8")
	err = dao.Database.DB.CreateMessage(ctx, &db.CreateMessageParams{
		NotifyType: db.MessagesNotifyTypeCommon,
		MsgType:    db.MessagesMsgType(model.MsgTypeText),
		MsgContent: param.MsgContent,
		MsgExtend:  msgExtend,
		AccountID:  sql.NullInt64{Int64: param.AccountID, Valid: true},
		RlyMsgID:   sql.NullInt64{Int64: rlyMsgID, Valid: rlyMsgID > 0},
		RelationID: param.RelationID,
	})
	if err != nil {
		global.Logger.Error(err.Error())
		return nil, errcode.ErrServer
	}
	//fmt.Println("SendMsg9")
	result, err := dao.Database.DB.CreateMessageReturn(ctx)
	global.Worker.SendTask(task.PublishMsg(reply.ParamMsgInfoWithRly{
		ParamMsgInfo: reply.ParamMsgInfo{
			ID:         rlyMsgID,
			NotifyType: string(db.MessagesNotifyTypeCommon),
			MsgType:    string(model.MsgTypeText),
			MsgContent: result.MsgContent,
			MsgExtend:  param.MsgExtend,
			AccountID:  param.AccountID,
			RelationID: param.RelationID,
			CreateAt:   result.CreateAt,
		},
		RlyMsg: rlyMsg,
	}))
	//fmt.Println("SendMsg10")
	return &client.HandleSendMsgRly{
		MsgID:    result.ID,
		CreateAt: result.CreateAt,
	}, nil
}

func (message) ReadMsg(ctx context.Context, params *model.HandleReadMsg) errcode.Err {
	ok, myErr := logic.ExistsSetting(ctx, params.ReaderID, params.RelationID)
	if myErr != nil {
		return myErr
	}
	if !ok {
		return errcodes.AuthPermissionsInsufficient
	}
	err := dao.Database.DB.UpdateMsgReads(ctx, params.ReaderID)
	if err != nil {
		global.Logger.Error(err.Error())
		return errcode.ErrServer
	}
	readMsgs, err := dao.Database.DB.UpdateMsgReadsReturn(ctx, params.RelationID)
	msgMap := make(map[int64][]int64)
	for _, readMsg := range readMsgs {
		msgID, accountID := readMsg.ID, readMsg.AccountID
		msgMap[msgID] = append(msgMap[msgID], accountID)
	}
	global.Worker.SendTask(task.ReadMsg(params.AccessToken, params.ReaderID, msgMap, params.MsgIDs))
	return nil
}
