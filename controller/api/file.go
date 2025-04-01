package api

import (
	"ChatRoom001/global"
	"ChatRoom001/logic"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/request"
	"ChatRoom001/pkg/gtype"
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/app"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type file struct {
}

func (file) PublishFile(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamPublishFile)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid)
		return
	}
	fmt.Println("params:", params.File)
	fileType, myErr := gtype.GetFileType(params.File)
	fmt.Println(fileType, myErr)
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		reply.Reply(errcode.ErrServer)
		return
	}
	if fileType != "img" && fileType != "png" && fileType != "jpg" {
		fileType = "file"
	}
	result, err := logic.Logics.File.PublishFile(ctx, model.PublishFile{
		File:       params.File,
		RelationID: params.RelationID,
		AccountID:  params.AccountID,
	})
	reply.Reply(err, result)
}

func (file) DeleteFile(ctx *gin.Context) {

}
