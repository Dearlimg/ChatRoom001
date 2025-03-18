package routers

import (
	"ChatRoom001/controller/api"
	"ChatRoom001/middlewares"
	"github.com/gin-gonic/gin"
)

type application struct{}

func (application) Init(routers *gin.RouterGroup) {
	r := routers.Group("application").Use(middlewares.MustAccount())
	{
		r.POST("create", api.Apis.Application.CreateApplication)
		r.DELETE("delete", api.Apis.Application.DeleteApplication)
		r.PUT("refuse", api.Apis.Application.RefuseApplication)
		r.PUT("accept", api.Apis.Application.AcceptApplication)
	}
}
