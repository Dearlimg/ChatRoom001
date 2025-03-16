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
			userGroup.POST("token", api.Apis.Account.GetAccountToken)
			userGroup.GET("infos/account", api.Apis.Account.GetAccountByUserID)
		}
		accountGroup := r.Group("").Use(middlewares.MustAccount())
		{
			accountGroup.PUT("update", api.Apis.Account.UpdateAccount)
		}
	}
}
