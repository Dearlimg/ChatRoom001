package api

import (
	"ChatRoom001/errcodes"
	"ChatRoom001/logic"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/request"
	"github.com/Dearlimg/Goutils/pkg/app"
	"github.com/gin-gonic/gin"
)

type user struct {
}

func (user) Login(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamLogin)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcodes.PasswordNotValid.WithDetails(err.Error()))
		return
	}
	result, err := logic.Logics.User.Login(ctx, params.Email, params.Password)
	reply.Reply(err, result)
}

func (user) Register(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamRegister)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcodes.PasswordNotValid.WithDetails(err.Error()))
		return
	}
	result, err := logic.Logics.User.Register(ctx, params.Email, params.Password, params.Code)
	reply.Reply(err, result)
}

func (user) Logout(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	if err := logic.Logics.User.Logout(ctx); err != nil {
		reply.Reply(err)
		return
	}
	reply.Reply(nil, gin.H{
		"msg": "登出成功",
	})
}

func (user) UpdateUserPassword(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateUserPassword)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcodes.PasswordNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.User.UpdateUserPassword(ctx, content.ID, params.Code, params.NewPassword)
	reply.Reply(err, nil)
}
