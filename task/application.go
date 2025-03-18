package task

import (
	"ChatRoom001/global"
	"ChatRoom001/model/chat"
)

func Application(accountID int64) func() {
	return func() {
		global.ChatMap.Send(accountID, chat.ServerApplication)
	}
}
