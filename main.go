package main

import (
	"ChatRoom001/global"
	"ChatRoom001/model/common"
	"ChatRoom001/routers/router"
	"ChatRoom001/setting"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	setting.Inits()
	if global.PublicSetting.Server.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("email", common.ValidatorEmail)
	}

	r, ws := router.NewRouter()
}
