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

type account struct {
}

func (account) CreateAccount(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamsCreateAccount)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok && content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Account.CreateAccount(ctx, content.ID, params.Name, global.PublicSetting.Rules.DefaultAvatarURL, params.Gender, params.Signature)
	reply.Reply(err, result)
}

func (account) GetAccountToken(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetAccountToken)
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)

	if !ok && content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	fmt.Println("GetAccountToken", content.ID, params.AccountID)
	result, err := logic.Logics.Account.GetAccountToken(ctx, content.ID, params.AccountID)
	reply.Reply(err, result)
}

func (account) GetAccountByUserID(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok && content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Account.GetAccountsByUserID(ctx, content.ID)
	reply.Reply(err, result)
}

func (account) UpdateAccount(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	param := new(request.ParamUpdateAccount)
	if err := ctx.ShouldBind(param); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	fmt.Println("Update Account", content, ok)

	if !ok && content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err1 := logic.Logics.Account.UpdateAccount(ctx, param.AccountID, param.Name, param.Gender, param.Signature)
	reply.Reply(err1, nil)
}

func (account) GetAccountsByName(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetAccountByName)
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	fmt.Println("Get Account By Name", params)
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok && content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Account.GetAccountsByName(ctx, content.ID, params.Name, limit, offset)
	reply.Reply(err, result)
}

func (account) GetAccountByID(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetAccountByID)
	//fmt.Println("GetAccountByID ", params)
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	fmt.Println("GetAccountByID ", content, ok, params)
	if !ok && content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Account.GetAccountByID(ctx, content.ID, params.AccountID)
	reply.Reply(err, result)
}

func (account) DeleteAccount(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	param := new(request.ParamDeleteAccount)
	if err := ctx.ShouldBindBodyWithJSON(param); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok && content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Account.DeleteAccount(ctx, content.ID, param.AccountID)
	reply.Reply(err, nil)
}
