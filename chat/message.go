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
	socketio "github.com/googollee/go-socket.io"
)

type message struct {
	server *socketio.Server //
}

var processedMessages = make(map[int64]bool)

func (message) SendMsg(ctx context.Context, param *model.HandleSendMsg) (*client.HandleSendMsgRly, errcode.Err) {
	// 检查消息是否已经处理过
	//if processedMessages[param.MsgID] {
	//	// 直接返回之前的处理结果，这里简单返回一个空结果
	//	return &client.HandleSendMsgRly{}, nil
	//}
	//
	//// 标记消息为已处理
	//processedMessages[param.MsgID] = true

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

	if param.RlyMsgID > 0 {
		rlyInfo, myErr := logic.GetMsgInfoByID(ctx, param.RlyMsgID)
		if myErr != nil {
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

		rlyMsg = &reply.ParamRlyMsg{
			MsgID:      rlyInfo.ID,
			MsgType:    string(rlyInfo.MsgType),
			MsgContent: rlyInfo.MsgContent,
			MsgExtend:  rlyMsgExtend,
			IsRevoked:  rlyInfo.IsRevoke,
		}
	}

	msgExtend, err := model.ExtendToJson(param.MsgExtend)
	if err != nil {
		global.Logger.Error(err.Error())
		return nil, errcode.ErrServer
	}

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
	result, err := dao.Database.DB.CreateMessageReturn(ctx)
	fmt.Println("\036[31mSend msg before !!\036[0m")
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
