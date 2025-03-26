package api

import (
	"ChatRoom001/errcodes"
	"ChatRoom001/logic"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/request"
	"github.com/Dearlimg/Goutils/pkg/app"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type message struct{}

func (message) CreateFileMsg(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamCreateFileMsg)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Message.CreateFileMsg(ctx, model.CreateFileMsg{
		AccountID:  content.ID,
		RelationID: params.RelationID,
		RlyMsgID:   params.RlyMsg,
		File:       params.File,
	})
	reply.Reply(err, result)
}
