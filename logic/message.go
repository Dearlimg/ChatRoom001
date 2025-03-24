package logic

import (
	"ChatRoom001/dao"
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/errcodes"
	"ChatRoom001/global"
	"context"
	"database/sql"
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/pkg/errors"
)

type message struct{}

func GetMsgInfoByID(ctx context.Context, msgID int64) (*db.GetMessageByIDRow, errcode.Err) {
	result, err := dao.Database.DB.GetMessageByID(ctx, msgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.MsgNotExists
		}
		global.Logger.Error(err.Error())
		return nil, errcode.ErrServer
	}
	return result, nil
}

//func (message) CreateFileMsg(ctx *gin.Context, params model.CreateFileMsg) (*reply.ParamCreateFileMsg, errcode.Err) {
//	ok, myErr := ExistsSetting(ctx, params.AccountID, params.RelationID)
//	if myErr != nil {
//		return nil, myErr
//	}
//	if !ok {
//		return nil, errcodes.AuthPermissionsInsufficient
//	}
//	fileInfo, myErr := Logics.File.Publish(ctx)
//}
