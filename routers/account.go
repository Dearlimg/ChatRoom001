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
			userGroup.GET("token", api.Apis.Account.GetAccountToken)
			userGroup.GET("infos/account", api.Apis.Account.GetAccountByUserID)
			userGroup.DELETE("delete", api.Apis.Account.DeleteAccount)
		}
		accountGroup := r.Group("").Use(middlewares.MustAccount())
		{
			accountGroup.PUT("update", api.Apis.Account.UpdateAccount)
			accountGroup.GET("info", api.Apis.Account.GetAccountByID)
			accountGroup.GET("infos/name", api.Apis.Account.GetAccountsByName)
		}
	}
}
