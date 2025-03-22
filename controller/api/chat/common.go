package chat

import (
	"ChatRoom001/dao"
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/chat"
	"ChatRoom001/model/common"
	"context"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	socketio "github.com/googollee/go-socket.io"
	"time"
)

func CheckConnCtxToken(v interface{}) (*model.Token, errcode.Err) {
	token, ok := v.(*model.Token)
	if !ok {
		return nil, errcodes.AuthenticationFailed
	}
	if token.PayLoad.ExpiredAt.Before(time.Now()) {
		return nil, errcodes.AuthOverTime
	}
	return token, nil
}

func CheckAuth(s socketio.Conn) (*model.Token, bool) {
	token, err := CheckConnCtxToken(s.Context())
	if err != nil {
		s.Emit(chat.ServerError, common.NewState(err))
		_ = s.Close()
		return nil, false
	}
	return token, true
}

func MustAccount(accessToken string) (*model.Token, errcode.Err) {
	payload, _, err := middlewares.ParseToken(accessToken)
	if err != nil {
		return nil, err
	}
	content := new(model.Content)
	if err := content.Unmarshal(payload.Content); err != nil {
		return nil, errcodes.AuthenticationFailed
	}
	ok, err1 := dao.Database.DB.ExistAccountByID(context.Background(), content.ID)
	if err1 != nil {
		global.Logger.Error(err1.Error())
		return nil, err
	}
	if !ok {
		return nil, errcodes.AccountNotFound
	}
	return &model.Token{
		AccessToken: accessToken,
		PayLoad:     payload,
		Content:     content,
	}, nil
}
