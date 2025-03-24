package logic

type file struct {
}

//func (file) PublishFile(ctx *gin.Context, params model.PublishFile) (model.PublishFileReply, errcode.Err) {
//	fileType, myErr := gtype.GetFileType(params.File)
//	if myErr != nil {
//		return model.PublishFileReply{}, myErr
//	}
//	if fileType == "file" {
//		if params.File.Size > global.PublicSetting.Rules.BiggestFileSize {
//			return model.PublishFileReply{}, myErr
//		}
//	} else {
//		fileType = "img"
//	}
//	input := new(oss.PutBucketCname)
//}

//// PublishFile 上传文件，传出 context 与 relationID，accountID，file(*multipart.FileHeader)，返回 model.PublishFileRe
//func (file) PublishFile(ctx *gin.Context, params model.PublishFile) (model.PublishFileReply, errcode.Err) {
//	// 文件类型验证逻辑保持不变
//	fileType, myErr := gtype.GetFileType(params.File)
//	if myErr != nil {
//		return model.PublishFileReply{}, errcode.ErrServer
//	}
//	if fileType == "file" {
//		if params.File.Size > global.PublicSetting.Rules.BiggestFileSize {
//			return model.PublishFileReply{}, errcodes.FileTooBig
//		}
//	} else {
//		fileType = "img"
//	}
//
//	// 阿里云OSS上传逻辑
//	fileHeader := params.File
//	file, err := fileHeader.Open()
//	if err != nil {
//		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
//		return model.PublishFileReply{}, errcode.ErrServer
//	}
//	defer file.Close()
//
//	// 生成唯一对象键（示例使用UUID+文件名）
//	fileExt := path.Ext(fileHeader.Filename)
//	key := fmt.Sprintf("uploads/%s%s", uuid.New().String(), fileExt)
//
//	// 获取Content-Type
//	contentType := fileHeader.Header.Get("Content-Type")
//	if contentType == "" {
//		contentType = mime.TypeByExtension(fileExt)
//	}
//
//	// 创建OSS PutObject选项
//	options := []oss.Option{
//		oss.ContentType(contentType),
//		oss.ObjectACL(oss.ACLPublicRead), // 根据需求设置ACL
//	}
//
//	// 执行上传
//	err = global.OSSBucket.PutObject(key, file, options...)
//	if err != nil {
//		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
//		return model.PublishFileReply{}, errcode.ErrServer
//	}
//
//	// 生成访问URL（假设使用公共读方式）
//	url := fmt.Sprintf("https://%s.%s/%s",
//		global.OSSConfig.BucketName,
//		global.OSSConfig.Endpoint,
//		key)
//
//	// 数据库操作保持不变
//	dao.Database.DB.CreateFile(ctx, &db.CreateFileParams{
//		FileName: fileHeader.Filename,
//		FileType: db.Filetype(fileType),
//		FileSize: fileHeader.Size,
//		Key:      key,
//		Url:      url,
//		RelationID: sql.NullInt64{
//			Int64: params.RelationID,
//			Valid: true,
//		},
//		AccountID: sql.NullInt64{
//			Int64: params.AccountID,
//			Valid: true,
//		},
//	})
//
//	r, err := dao.Database.DB.CreateFileReturn(ctx)
//
//	if err != nil {
//		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
//		return model.PublishFileReply{}, errcode.ErrServer
//	}
//
//	return model.PublishFileReply{
//		ID:       r.ID,
//		FileType: fileType,
//		FileSize: r.FileSize,
//		Url:      r.Url,
//		CreateAt: r.CreateAt,
//	}, nil
//}
