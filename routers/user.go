package routers

import (
	"ChatRoom001/controller/api"
	"github.com/gin-gonic/gin"
)

type user struct {
}

func (user) Init(router *gin.RouterGroup) {
	r := router.Group("user")
	{
		r.POST("register", api.Apis.User.Register)
		r.POST("/login", api.Apis.User.Login)
		r.DELETE("/logout", api.Apis.User.Logout)
	}
	r.DELETE("deleteUser")
}
