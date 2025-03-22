package routers

import (
	"ChatRoom001/controller/api"
	"ChatRoom001/middlewares"
	"github.com/gin-gonic/gin"
)

type message struct{}

func (message) Init(router *gin.RouterGroup) {
	r := router.Group("message").Use(middlewares.MustAccount())
	{
		r.POST("file", api.Apis.Message.CreateFileMsg)
	}
}
