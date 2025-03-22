package routers

import (
	"ChatRoom001/controller/api"
	chat2 "ChatRoom001/model/chat"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type ws struct {
}

func (ws) Init(router *gin.Engine) *socketio.Server {
	server := socketio.NewServer(nil)
	{
		server.OnConnect("/", api.Apis.Chat.Handle.OnConnect)
		server.OnError("/", api.Apis.Chat.Handle.OnError)
	}
	chatHande(server)
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	return server
}

func chatHande(server *socketio.Server) {
	namespace := "/chat"
	server.OnEvent(namespace, chat2.ClientSendMsg, api.Apis.Chat.Message.SendMsg)
	server.OnEvent(namespace, chat2.ClientReadMsg, api.Apis.Chat.Message.ReadMsg)
	server.OnEvent(namespace, chat2.ClientAuth, api.Apis.Chat.Handle.Auth) // 账户登录
	server.OnEvent(namespace, chat2.ClientTest, api.Apis.Chat.Handle.Test)
}
