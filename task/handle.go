package task

import (
	"ChatRoom001/global"
	"ChatRoom001/model/chat"
	"ChatRoom001/model/chat/server"
	"github.com/Dearlimg/Goutils/pkg/utils"
)

func AccountLogin(accessToken, address string, accountID int64) func() {
	return func() {
		global.ChatMap.Send(accountID, chat.ServerAccountLogin, server.AccessLogin{
			Address: address,
			EnToken: utils.EncodeMD5(accessToken),
		})
	}
}
