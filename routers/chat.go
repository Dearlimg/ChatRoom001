package routers

import (
	"ChatRoom001/controller/api"
	chat2 "ChatRoom001/model/chat"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"net/http"
)

type ws struct {
}

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func (ws) Init(router *gin.Engine) *socketio.Server {
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})
	{
		server.OnConnect("/", api.Apis.Chat.Handle.OnConnect)
		server.OnError("/", api.Apis.Chat.Handle.OnError)
		server.OnDisconnect("/", api.Apis.Chat.Handle.OnDisconnect)
	}
	chatHande(server)
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	return server
}

func chatHande(server *socketio.Server) {
	namespace := "/"
	server.OnEvent(namespace, chat2.ClientSendMsg, api.Apis.Chat.Message.SendMsg)
	server.OnEvent(namespace, chat2.ClientReadMsg, api.Apis.Chat.Message.ReadMsg)
	server.OnEvent(namespace, chat2.ClientAuth, api.Apis.Chat.Handle.Auth) // 账户登录
	server.OnEvent(namespace, chat2.ClientTest, api.Apis.Chat.Handle.Test)
}
