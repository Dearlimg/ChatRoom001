package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/testdata/enums/api"
)

type user struct {
}

func (user) Init(router *gin.RouterGroup) {
	r := router.Group("user")
	{
		r.POST("register", api.Apis.User.register)
	}
}
