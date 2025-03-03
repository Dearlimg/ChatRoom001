package api

import (
	"ChatRoom001/errcodes"
	"ChatRoom001/logic"
	"ChatRoom001/model/request"
	"github.com/Dearlimg/Goutils/pkg/app"
	"github.com/gin-gonic/gin"
)

type user struct {
}


func (user) Register(ctx *gin.Context) {
	reply:=app.NewResponse(ctx)
	params:=new(request.ParamRegister)
	if err:=ctx.ShouldBind(params);err!=nil{
		reply.Reply(errcodes.PasswordNotValid.WithDetails(err.Error()))
		return
	}

	result,err:=logic.

}