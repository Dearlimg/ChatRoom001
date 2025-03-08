package routers

import (
	"ChatRoom001/controller/api"
	"ChatRoom001/middlewares"
	"github.com/gin-gonic/gin"
)

type account struct {
}

func (account) Init(routers *gin.RouterGroup) {
	r := routers.Group("account")
	{
		userGroup := r.Group("").Use(middlewares.MustUser())
		{
			userGroup.POST("create", api.Apis.Account.CreateAccount)
		}
	}
}
