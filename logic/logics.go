package logic

type Logic struct {
	Email       email
	User        user
	Account     account
	Application application
}

var Logics = new(Logic)
