package routers

type routers struct {
	User    user
	Account account
	Chat    chat
}

var Routers = new(routers)
