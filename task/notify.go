package task

import (
	"ChatRoom001/dao"
	"ChatRoom001/global"
	"ChatRoom001/model"
	"ChatRoom001/model/chat"
	"ChatRoom001/model/chat/server"
	"github.com/Dearlimg/Goutils/pkg/utils"
)

func CreateNotify(accessToken string, accountID, relationID int64, msgContent string, msgExtend *model.MsgExtend) func() {
	ctx, cancel := global.DefaultContextWithTimeout()
	defer cancel()
	member, err := dao.Database.DB.GetGroupMembers(ctx, relationID)
	if err != nil {
		global.Logger.Error(err.Error())
	}
	return func() {
		global.ChatMap.SendMany(member, chat.ServerCreateNotify, server.CreateNotify{
			EnToken:    utils.EncodeMD5(accessToken),
			AccountID:  accountID,
			RelationID: relationID,
			MsgContent: msgContent,
			MsgExtend:  msgExtend,
		})
	}
}

func UpdateNotify(accessToken string, accountID, relationID int64, msgContent string, msgExtend *model.MsgExtend) func() {
	ctx, cancel := global.DefaultContextWithTimeout()
	defer cancel()
	members, err := dao.Database.DB.GetGroupMembers(ctx, relationID)
	if err != nil {
		global.Logger.Error(err.Error())
	}
	return func() {
		global.ChatMap.SendMany(members, chat.ServerUpdateNotify, server.CreateNotify{
			EnToken:    utils.EncodeMD5(accessToken),
			AccountID:  accountID,
			RelationID: relationID,
			MsgContent: msgContent,
			MsgExtend:  msgExtend,
		})
	}
}
