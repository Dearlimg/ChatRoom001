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
			userGroup.PUT("update", api.Apis.Account.UpdateAccount)
			userGroup.DELETE("delete", api.Apis.Account.DeleteAccount)
			userGroup.GET("infos/ID", api.Apis.Account.GetAccountByID)
		}
		accountGroup := r.Group("").Use(middlewares.MustAccount())
		{
			accountGroup.GET("infos/name", api.Apis.Account.GetAccountsByName)
		}
	}
}
