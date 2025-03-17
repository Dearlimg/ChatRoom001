package setting

import (
	"ChatRoom001/global"
	"ChatRoom001/manager"
)

type chat struct{}

func (chat) Init() {
	global.ChatMap = manager.NewChatMap()
}
