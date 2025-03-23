package routers

import (
	"ChatRoom001/controller/api"
	"ChatRoom001/middlewares"
	"github.com/gin-gonic/gin"
)

type file struct{}

func (file) Init(routers *gin.RouterGroup) {
	r := routers.Group("file", middlewares.MustAccount())
	{
		r.POST("publish", api.Apis.File.PublishFile)
	}
}
