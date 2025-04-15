package task

import (
	"ChatRoom001/dao"
	"ChatRoom001/global"
	"ChatRoom001/model/chat"
	"ChatRoom001/model/chat/server"
	"github.com/Dearlimg/Goutils/pkg/utils"
)

func TransferGroup(accessToken string, accountID, relationID int64) func() {
	ctx, cancel := global.DefaultContextWithTimeout()
	defer cancel()
	// 获取群中所有成员的ID
	members, err := dao.Database.Redis.GetAllAccountsByRelationID(ctx, relationID)
	if err != nil {
		global.Logger.Error(err.Error())
	}
	return func() {
		global.ChatMap.SendMany(members, chat.ServerGroupTransferred, server.TransferGroup{
			EnToken:   utils.EncodeMD5(accessToken),
			AccountID: accountID,
		})
	}
}
