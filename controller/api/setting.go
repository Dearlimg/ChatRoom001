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

type setting struct{}

func (setting) GetFriends(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Setting.GetFriends(ctx, content.ID)
	reply.Reply(err, result)
}

func (setting) GetFriendsByName(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetFriendsByName)
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Setting.GetFriendsByName(ctx, content.ID, params.Name, limit, offset)
	reply.Reply(err, result)
}

func (setting) GetShows(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Setting.GetShows(ctx, content.ID)
	reply.Reply(err, result)
}

func (setting) GetPins(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Setting.GetPins(ctx, content.ID)
	reply.Reply(err, result)
}

func (setting) UpdateSettingPin(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamUpdateSettingPin)
	if err := ctx.ShouldBind(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Setting.UpdatePin(ctx, content.ID, params.RelationID, *params.IsPin)
	rly.Reply(err, nil)
}

func (setting) UpdateNickName(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := &request.ParamUpdateNickName{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Setting.UpdateNickName(ctx, content.ID, params.RelationID, params.NickName)
	reply.Reply(err, nil)
}

func (setting) UpdateSettingShow(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := &request.ParamUpdateSettingShow{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Setting.UpdateShow(ctx, content.ID, params.RelationID, *params.IsShow)
	reply.Reply(err, nil)
}

func (setting) UpdateSettingDisturb(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := &request.ParamUpdateSettingDisturb{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Setting.UpdateDisaturb(ctx, content.ID, params.RelationID, *params.IsNotDisturb)
	reply.Reply(err, nil)
}
