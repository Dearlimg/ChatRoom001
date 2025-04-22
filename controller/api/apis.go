package api

import "ChatRoom001/controller/api/chat"

type apis struct {
	User        user
	Account     account
	Chat        chat.Group
	Email       email
	Group       group
	Application application
	Message     message
	File        file
	Setting     setting
	Notify      notify
}

var Apis = new(apis)
