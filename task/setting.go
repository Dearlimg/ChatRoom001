package task

import (
	"ChatRoom001/global"
	"ChatRoom001/model/chat"
	"ChatRoom001/model/chat/server"
	"github.com/Dearlimg/Goutils/pkg/utils"
)

func UpdateNickName(accessToken string, accountID, relationID int64, nickName string) func() {
	return func() {
		global.ChatMap.Send(accountID, chat.ServerUpdateNickName, server.UpdateNickName{
			EnToken:    utils.EncodeMD5(accessToken),
			RelationID: relationID,
			NickName:   nickName,
		})
	}
}
