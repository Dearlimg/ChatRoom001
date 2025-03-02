package chat

import (
	"ChatRoom001/logic"
	"ChatRoom001/model/request"
	"github.com/Dearlimg/Goutils/pkg/app"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type email struct {
}

func (email) ExistEmail(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := &request.ParamExistEmail{}
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	result, err := logic.Logics.Email.ExistEmail(ctx, params.Email)
	reply.Reply(err, result)
}

func (email) SendMark(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := &request.ParamSendEmail{}
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
	}
	err := logic.Logics.Email.SendMark(params.Email)
}
