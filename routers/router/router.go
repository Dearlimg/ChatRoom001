package router

import (
	"ChatRoom001/global"
	"ChatRoom001/middlewares"
	"ChatRoom001/routers"
	"github.com/Dearlimg/Goutils/pkg/app"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func NewRouter() (*gin.Engine, *socketio.Server) {
	r := gin.New()
	r.Use(middlewares.Cors(), middlewares.GinLogger(), middlewares.Recovery(true))
	root := r.Group("api", middlewares.LogBody(), middlewares.PasetoAuth())
	{
		root.GET("swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
		root.GET("ping", func(ctx *gin.Context) {
			reply := app.NewResponse(ctx)
			global.Logger.Info("ping", middlewares.ErrLogMsg(ctx)...)
			reply.Reply(nil, "pong")
		})
		root.GET("/man", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"msg0": "haha",
				"msg1": "what can i say",
				"msg2": "manBa out",
				"msg3": " 孩子们,这并不好笑.",
			})
		})
		rg := routers.Routers
		rg.Account.Init(root) //finish
		rg.User.Init(root)    //need modify
		rg.Email.Init(root)   //finish
		rg.Group.Init(root)
		rg.Application.Init(root) //finish
		rg.File.Init(root)        //finish
		rg.Message.Init(root)     //finish
		rg.Setting.Init(root)     //need delete
	}
	return r, routers.Routers.Chat.Init(r)
}
