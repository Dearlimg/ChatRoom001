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
}

var Routers = new(routers)
