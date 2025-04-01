package logic

type Logic struct {
	Email       email
	User        user
	Account     account
	Application application
	File        file
	Message     message
	Setting     setting
}

var Logics = new(Logic)
