package routers

type routers struct {
	User    user
	Account account
	Chat    ws
}

var Routers = new(routers)
