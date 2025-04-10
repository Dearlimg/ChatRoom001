package api

import (
	"ChatRoom001/errcodes"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/request"
	"github.com/Dearlimg/Goutils/pkg/app"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type group struct{}

func (group) CreateGroup(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamCreateGroup)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}

	//err := logic.Logics.Group.CreateGroup(ctx, content.ID, params.Name, params.Description)
	//if err != nil {
	//	reply.Reply(err)
	//	return
	//}
}
