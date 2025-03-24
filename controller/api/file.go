package api

type file struct {
}

//func (file) PublishFile(ctx *gin.Context) {
//	reply := app.NewResponse(ctx)
//	params := new(request.ParamPublishFile)
//	if err := ctx.ShouldBind(params); err != nil {
//		reply.Reply(errcode.ErrParamsNotValid)
//		return
//	}
//	fileType, myErr := gtype.GetFileType(params.File)
//	if myErr != nil {
//		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
//		reply.Reply(errcode.ErrServer)
//		return
//	}
//	if fileType != "img" && fileType != "png" && fileType != "jpg" {
//		fileType = "file"
//	}
//	result, err := logic.Logics.File.PublishFile(ctx, model.PublishFile{
//		File:       params.File,
//		RelationID: params.RelationID,
//		AccountID:  params.AccountID,
//	})
//	reply.Reply(err, result)
//}
