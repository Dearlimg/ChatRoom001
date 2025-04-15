package api

import (
	"ChatRoom001/errcodes"
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
