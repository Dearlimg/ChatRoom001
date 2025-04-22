package routers

type routers struct {
	User        user
	Account     account
	Chat        ws
	Email       email
	Group       group
	Application application
	Message     message
	File        file
	Setting     setting
	Notify      notify
}

var Routers = new(routers)
