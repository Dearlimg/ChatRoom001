package api

import (
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/logic"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/reply"
	"ChatRoom001/model/request"
	"github.com/Dearlimg/Goutils/pkg/app"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type group struct{}

func (group) CreateGroup(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamCreateGroup)
	if err := ctx.ShouldBindJSON(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}

	relationID, err := logic.Logics.Group.CreateGroup(ctx, content.ID, params.Name, params.Description)
	if err != nil {
		rly.Reply(err)
		return
	}
	result, err := logic.Logics.File.UploadGroupAvatar(ctx, nil, content.ID, relationID)
	rly.Reply(err, reply.ParamCreateGroup{
		Name:        params.Name,
		AccountID:   content.ID,
		RelationID:  relationID,
		Description: params.Description,
		Avatar:      result.URL,
	})
}

func (group) TransferGroup(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamTransferGroup)
	if err := ctx.ShouldBindJSON(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Group.TransferGroup(ctx, content.ID, params.RelationID, params.ToAccountID)
	rly.Reply(err)
}

func (group) DissolveGroup(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamDissolveGroup)
	if err := ctx.ShouldBindJSON(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Group.DissolveGroup(ctx, content.ID, params.RelationID)
	rly.Reply(err)
}

func (group) UpdateGroup(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamUpdateGroup)
	if err := ctx.ShouldBind(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Group.UpdateGroup(ctx, content.ID, params.RelationID, params.Name, params.Description)
	if err != nil {
		rly.Reply(err, result)
	}
	avatar, err := logic.Logics.File.UploadGroupAvatar(ctx, params.Avatar, content.ID, params.RelationID)
	result.Avatar = avatar.URL
	rly.Reply(err, result)
}

func (group) InviteAccount(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamInviteAccount)
	if err := ctx.ShouldBindJSON(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Group.InviteAccount(ctx, content.ID, params.RelationID, params.AccountID)
	rly.Reply(err, result)
}

func (group) GetGroupList(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Group.GetGroupList(ctx, content.ID)
	rly.Reply(err, result)
}

func (group) QuitGroup(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamQuitGroup)
	if err := ctx.ShouldBindJSON(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Group.QuitGroup(ctx, content.ID, params.RelationID)
	rly.Reply(err)
}

func (group) GetGroupsByName(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamGetGroupsByName)
	if err := ctx.ShouldBindJSON(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Group.GetGroupsByName(ctx, content.ID, params.Name, limit, offset)
	rly.Reply(err, result)
}

func (group) GetGroupMembers(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamGetGroupMembers)
	if err := ctx.ShouldBindQuery(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Group.GetGroupMembers(ctx, content.ID, params.RelationID, limit, offset)
	rly.Reply(err, result)
}
