package api

import "ChatRoom001/controller/api/chat"

type apis struct {
	User    user
	Account account
	Chat    chat.Group
}

var Apis = new(apis)
