package api

import (
	"ChatRoom001/errcodes"
	"ChatRoom001/logic"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/request"
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/app"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type notify struct{}

func (notify) CreateNotify(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamCreateNotify)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Notify.CreateNotify(ctx, content.ID, params)
	reply.Reply(err, result)
}

func (notify) UpdateNotify(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateNotify)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Notify.UpdateNotify(ctx, content.ID, params)
	reply.Reply(err, nil)
}

func (notify) GetNotifyByID(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetNotifyByID)
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	fmt.Println("GetNotifyByID content:", content, params)
	data, err := logic.Logics.Notify.GetNotifyByID(ctx, content.ID, params)
	reply.Reply(err, data)
}

func (notify) DeleteNotify(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamDeleteNotify)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Notify.DeleteNotify(ctx, content.ID, params)
	reply.Reply(err, nil)
}
