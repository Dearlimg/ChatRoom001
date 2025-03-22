package chat

import (
	"ChatRoom001/global"
	"ChatRoom001/model/chat/client"
	"ChatRoom001/model/common"
	"ChatRoom001/pkg/mq/consumer"
	"ChatRoom001/task"
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"time"
)

type handle struct {
}

const AuthLimitTimeout = 10 * time.Second

func (handle) OnConnect(s socketio.Conn) error {
	log.Println("connected", s.RemoteAddr().String(), s.ID()) // s.RemoteAddr()获取客户端的 IP 地址和端口号信息。
	time.AfterFunc(AuthLimitTimeout, func() {
		if !global.ChatMap.HasSID(s.ID()) {
			global.Logger.Info(fmt.Sprintln("auth failed:", s.RemoteAddr().String(), s.ID()))
			_ = s.Close()
		}
	})
	return nil
}

func (handle) OnError(s socketio.Conn, err error) {
	log.Println("on error:", err)
	if s == nil {
		return
	}
	global.ChatMap.Leave(s)
	log.Println("disconnected:", s.RemoteAddr().String(), s.ID())
	_ = s.Close()
}

func (handle) Auth(s socketio.Conn, accessToken string) string {
	token, myErr := MustAccount(accessToken)
	if myErr != nil {
		return common.NewState(myErr).MustJson()
	}
	s.SetContext(token)
	global.ChatMap.Link(s, token.Content.ID)
	global.Worker.SendTask(task.AccountLogin(accessToken, s.RemoteAddr().String(), token.Content.ID))
	log.Println("auth accept:", s.RemoteAddr().String())
	go consumer.StartConsumer(token.Content.ID)
	return common.NewState(nil).MustJson()
}

func (handle) Test(s socketio.Conn, msg string) string {
	_, ok := CheckAuth(s)
	if !ok {
		return ""
	}
	param := new(client.TestParams)
	log.Println(msg)
	if err := common.Decode(msg, param); err != nil {
		return common.NewState(errcode.ErrParamsNotValid.WithDetails(err.Error())).MustJson()
	}
	result := common.NewState(nil, client.TestRly{
		Name:    param.Name,
		Age:     param.Age,
		Address: s.RemoteAddr().String(),
		ID:      s.ID(),
	}).MustJson()
	s.Emit("test", "test")
	return result
}
