package logic

import (
	"ChatRoom001/dao"
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"ChatRoom001/middlewares"
	"ChatRoom001/model"
	"ChatRoom001/model/reply"
	"ChatRoom001/pkg/gtype"
	"database/sql"
	"fmt"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	upload "github.com/Dearlimg/Goutils/pkg/upload/obs"
	"github.com/Dearlimg/Goutils/pkg/upload/obs/ali_cloud"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"mime/multipart"
)

type file struct {
}

func initOSS(avatarType string) upload.OSS {
	if avatarType == ali_cloud.AccountAvatarType {
		return ali_cloud.Init(ali_cloud.Config{
			Location:         global.PrivateSetting.HuaWeiOBS.Location,
			BucketName:       global.PrivateSetting.HuaWeiOBS.BucketName,
			BucketUrl:        global.PrivateSetting.HuaWeiOBS.BucketUrl,
			Endpoint:         global.PrivateSetting.HuaWeiOBS.Endpoint,
			BasePath:         global.PrivateSetting.HuaWeiOBS.BasePath,
			AvatarType:       ali_cloud.AccountAvatarType,
			AccountAvatarUrl: global.PrivateSetting.HuaWeiOBS.AccountAvatarUrl,
			GroupAvatarUrl:   global.PrivateSetting.HuaWeiOBS.GroupAvatarUrl,
		})
	} else if avatarType == ali_cloud.GroupAvatarType {
		return ali_cloud.Init(ali_cloud.Config{
			Location:         global.PrivateSetting.HuaWeiOBS.Location,
			BucketName:       global.PrivateSetting.HuaWeiOBS.BucketName,
			BucketUrl:        global.PrivateSetting.HuaWeiOBS.BucketUrl,
			Endpoint:         global.PrivateSetting.HuaWeiOBS.Endpoint,
			BasePath:         global.PrivateSetting.HuaWeiOBS.BasePath,
			AvatarType:       ali_cloud.AccountAvatarType,
			AccountAvatarUrl: global.PrivateSetting.HuaWeiOBS.AccountAvatarUrl,
			GroupAvatarUrl:   global.PrivateSetting.HuaWeiOBS.GroupAvatarUrl,
		})
	}
	return global.OSS
}

func (file) PublishFile(ctx *gin.Context, params model.PublishFile) (model.PublishFileReply, errcode.Err) {
	fileType, myErr := gtype.GetFileType(params.File)
	if myErr != nil {
		return model.PublishFileReply{}, myErr
	}
	if fileType == "file" {
		if params.File.Size > global.PublicSetting.Rules.BiggestFileSize {
			return model.PublishFileReply{}, myErr
		}
		//}
	} else {
		fileType = "img"
	}
	url, key, err := global.OSS.UploadFile(params.File)
	if err != nil {
		fmt.Println("-------------------------------------------------------------------------------")
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return model.PublishFileReply{}, errcode.ErrServer
	}
	err = dao.Database.DB.CreateFile(ctx, &db.CreateFileParams{
		FileName: params.File.Filename,
		FileType: db.FilesFileType(fileType),
		FileSize: params.File.Size,
		Key:      key,
		Url:      url,
		RelationID: sql.NullInt64{
			Int64: params.RelationID,
			Valid: true,
		},
		AccountID: sql.NullInt64{
			Int64: params.AccountID,
			Valid: true,
		},
	})
	if err != nil {
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return model.PublishFileReply{}, errcode.ErrServer
	}
	r, err := dao.Database.DB.GetCreateFile(ctx, key)
	if err != nil {
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return model.PublishFileReply{}, errcode.ErrServer
	}
	return model.PublishFileReply{
		ID:       r.ID,
		FileType: fileType,
		FileSize: r.FileSize,
		Url:      r.Url,
		CreateAt: r.CreateAt,
	}, nil
}

func (file) UploadGroupAvatar(ctx *gin.Context, file *multipart.FileHeader, accountID, relationID int64) (*reply.ParamUploadAvatar, errcode.Err) {
	ok, err := dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return &reply.ParamUploadAvatar{URL: ""}, errcode.ErrServer
	}
	if !ok {
		return &reply.ParamUploadAvatar{URL: ""}, errcodes.NotGroupMember
	}
	oss := initOSS(ali_cloud.GroupAvatarType)
	var url, key string
	if file != nil {
		url, key, err = oss.UploadFile(file)
		if err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return &reply.ParamUploadAvatar{URL: ""}, errcode.ErrServer
		}
	}
	if file == nil {
		url = global.PublicSetting.Rules.DefaultAvatarURL
	}

	//filetype, myErr := gtype.GetFileType(file)
	//if myErr != nil {
	//	return &reply.ParamUploadAvatar{URL: ""}, errcode.ErrServer
	//}

	err = dao.Database.DB.UploadGroupAvatarWithTx(ctx, db.CreateFileParams{
		FileName: "groupAvatar",
		FileType: "",
		//FileType: db.FilesFileType(filetype),
		//FileSize:   0,
		FileSize:   file.Size,
		Key:        key,
		Url:        url,
		RelationID: sql.NullInt64{Int64: relationID, Valid: true},
		AccountID:  sql.NullInt64{Int64: accountID, Valid: true},
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return &reply.ParamUploadAvatar{URL: ""}, errcode.ErrServer
	}
	if file == nil {
		return &reply.ParamUploadAvatar{URL: global.PublicSetting.Rules.DefaultAvatarURL}, nil
	}
	return &reply.ParamUploadAvatar{URL: url}, nil
}

func (file) DeleteFile(ctx *gin.Context, fileID int64) errcode.Err {
	key, myErr := dao.Database.DB.GetFileKeyByID(ctx, fileID)
	if myErr != nil {
		if errors.Is(myErr, sql.ErrNoRows) {
			return errcodes.FileNotExist
		}
		global.Logger.Error(myErr.Error())
		return errcode.ErrServer
	}
	if key != "" {
		_, err := global.OSS.DeleteFile(key)
		if err != nil {
			global.Logger.Error(err.Error())
			return errcodes.FileDeleteFailed
		}
	}
	err := dao.Database.DB.DeleteFileByID(ctx, fileID)
	if err != nil {
		global.Logger.Error(err.Error())
		return errcode.ErrServer
	}
	return nil
}

func (file) UploadAccountAvatar(ctx *gin.Context, accountID int64, fileInfo *multipart.FileHeader) (*reply.ParamUploadAvatar, errcode.Err) {
	relationID, err := dao.Database.DB.GetRelationIDByAccountID(ctx, &db.GetRelationIDByAccountIDParams{
		Account1ID: sql.NullInt64{
			Int64: accountID,
			Valid: true,
		},
		Account2ID: sql.NullInt64{
			Int64: accountID,
			Valid: true,
		},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.RelationNotExists
		}
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	exists, err := dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if !exists {
		return nil, errcodes.AuthenticationFailed
	}
	var url string
	if fileInfo != nil {
		oss := initOSS(ali_cloud.AccountAvatarType)
		url, _, err = oss.UploadFile(fileInfo)
		if err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return nil, errcode.ErrServer
		}
	}
	err = dao.Database.DB.UpdateAccountAvatar(ctx, &db.UpdateAccountAvatarParams{
		Avatar: url,
		ID:     accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if fileInfo == nil {
		return &reply.ParamUploadAvatar{URL: global.PublicSetting.Rules.DefaultAvatarURL}, nil
	}
	return &reply.ParamUploadAvatar{URL: url}, nil
}

func (file) GetFileDetailsByID(ctx *gin.Context, fileID int64) (*reply.ParamFile, errcode.Err) {
	result, err := dao.Database.DB.GetFileDetailsByID(ctx, fileID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return &reply.ParamFile{
		FileID:    result.ID,
		FileName:  result.FileName,
		FileType:  string(result.FileType),
		FileSize:  result.FileSize,
		Url:       result.Url,
		AccountID: result.AccountID.Int64,
		CreateAt:  result.CreateAt,
	}, nil
}

func (file) GetRelationFile(ctx *gin.Context, relationID int64) (*reply.ParamGetRelationFile, errcode.Err) {
	list, err := dao.Database.DB.GetFileByRelationID(ctx, sql.NullInt64{Int64: relationID, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.FileNotExist
		}
	}
	data := make([]*reply.ParamFile, len(list))
	for i, f := range list {
		data[i] = &reply.ParamFile{
			FileID:    f.ID,
			FileName:  f.FileName,
			FileType:  string(f.FileType),
			FileSize:  f.FileSize,
			Url:       f.Url,
			AccountID: f.AccountID.Int64,
			CreateAt:  f.CreateAt,
		}
	}
	return &reply.ParamGetRelationFile{FileList: data}, nil
}
