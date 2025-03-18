package routers

type routers struct {
	User        user
	Account     account
	Chat        ws
	Email       email
	Group       group
	Application application
}

var Routers = new(routers)
