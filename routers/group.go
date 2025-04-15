package routers

import (
	"ChatRoom001/controller/api"
	"ChatRoom001/middlewares"
	"github.com/gin-gonic/gin"
)

type group struct {
}

func (group) Init(router *gin.RouterGroup) {
	r := router.Group("group").Use(middlewares.MustAccount())
	{
		r.POST("create", api.Apis.Group.CreateGroup)
		r.POST("transfer", api.Apis.Group.TransferGroup)
	}
}
