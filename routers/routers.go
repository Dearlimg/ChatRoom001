package routers

type routers struct {
	User    user
	Account account
	Chat    ws
	Email   email
}

var Routers = new(routers)
