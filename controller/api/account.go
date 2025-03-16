package api

import (
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/logic"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/request"
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
	param := new(request.ParamGetAccountToken)
	if err := ctx.ShouldBind(param); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)

	if !ok && content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}

	result, err := logic.Logics.Account.GetAccountToken(ctx, content.ID, param.AccountID)
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
	if !ok && content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Account.UpdateAccount(ctx, content.ID, param.Name, param.Gender, param.Signature)
	reply.Reply(err, nil)
}
