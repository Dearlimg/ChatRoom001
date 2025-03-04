package routers

import (
	//"ChatRoom001/controller/api"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type ws struct {
}

func (ws) Init(routers *gin.Engine) *socketio.Server {
	server := socketio.NewServer(nil)
	{

	}
	return server
}
