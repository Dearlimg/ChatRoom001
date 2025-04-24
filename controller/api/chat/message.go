package chat

import (
	"ChatRoom001/chat"
	"ChatRoom001/global"
	"ChatRoom001/model"
	"ChatRoom001/model/chat/client"
	"ChatRoom001/model/common"
	"encoding/base64"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	socketio "github.com/googollee/go-socket.io"
)

type message struct {
	server *socketio.Server
}

func (m message) SendMsg(s socketio.Conn, msg string) string {
	token, ok := CheckAuth(s)
	if !ok {
		return ""
	}
	param := new(client.HandleSendMsgParams)
	if err := common.Decode(msg, param); err != nil {
		return common.NewState(errcode.ErrParamsNotValid.WithDetails(err.Error())).MustJson()
	}
	temp, err := base64.StdEncoding.DecodeString(param.MsgContent)
	if err != nil {
		return err.Error()
	}
	param.MsgContent = string(temp)
	ctx, cancel := global.DefaultContextWithTimeout()
	defer cancel()
	//var count int64
	//count = 0
	result, myErr := chat.Group.Message.SendMsg(ctx, &model.HandleSendMsg{
		//MsgID:       count,
		AccessToken: token.AccessToken,
		RelationID:  param.RelationID,
		AccountID:   token.Content.ID,
		MsgContent:  param.MsgContent,
		MsgExtend:   param.MsgExtend,
		RlyMsgID:    param.RlyMsgID,
	})
	//count++
	return common.NewState(myErr, result).MustJson()
}

func (message) ReadMsg(s socketio.Conn, msg string) string {
	token, ok := CheckAuth(s)
	if !ok {
		return ""
	}
	param := new(client.HandleReadMsgParams)
	if err := common.Decode(msg, param); err != nil {
		return common.NewState(errcode.ErrParamsNotValid.WithDetails(err.Error())).MustJson()
	}
	ctx, cancel := global.DefaultContextWithTimeout()
	defer cancel()
	myErr := chat.Group.Message.ReadMsg(ctx, &model.HandleReadMsg{
		AccessToken: token.AccessToken,
		RelationID:  param.RelationID,
		MsgIDs:      param.MsgIDs,
		ReaderID:    token.Content.ID,
	})
	return common.NewState(myErr).MustJson()
}
