package api

import (
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/logic"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/request"
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/app"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type application struct {
}

func (application) CreateApplication(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	param := new(request.ParamCreateApplication)
	if err := ctx.ShouldBindJSON(param); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	context, ok := middlewares.GetTokenContent(ctx)
	if !ok {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	//fmt.Println("API Create Application", context, param.AccountID, param.ApplicationMsg)
	err := logic.Logics.Application.CreateApplication(ctx, context.ID, param.AccountID, param.ApplicationMsg)
	reply.Reply(err)
}

func (application) DeleteApplication(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	param := new(request.ParamDeleteApplication)
	if err := ctx.ShouldBindJSON(param); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	context, ok := middlewares.GetTokenContent(ctx)
	if !ok {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Application.DeleteApplication(ctx, context.ID, param.AccountID)
	reply.Reply(err)
}

func (application) RefuseApplication(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	param := new(request.ParamRefuseApplication)
	if err := ctx.ShouldBindJSON(param); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Application.RefuseApplication(ctx, param.AccountID, content.ID, param.RefuseMsg)
	reply.Reply(err)
}

func (application) AcceptApplication(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	param := new(request.ParamAcceptApplication)
	if err := ctx.ShouldBindJSON(param); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok && content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Application.AcceptApplication(ctx, content.ID, param.AccountID)
	fmt.Println(err)
}

func (application) ListApplications(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Application.ListApplications(ctx, content.ID, limit, offset)
	reply.Reply(err, result)

}
