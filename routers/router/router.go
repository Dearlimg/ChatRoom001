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
				"msg2": "manba out",
			})
		})
		rg := routers.Routers
		rg.Account.Init(root)
		rg.User.Init(root)
		rg.Email.Init(root)
		rg.Group.Init(root)
		rg.Application.Init(root)
		rg.File.Init(root)
		rg.Message.Init(root)
	}
	return r, routers.Routers.Chat.Init(r)
}
