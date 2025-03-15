package setting

type group struct {
	Config     config
	Dao        database
	Logger     log
	EmailMark  mark
	Worker     worker
	TokenMaker tokenMaker
	GenerateID generateID
}

var Group group

func Inits() {
	Group.Config.Init()
	Group.Dao.Init()
	Group.Logger.Init()
	Group.EmailMark.Init()
	Group.Worker.Init()
	Group.TokenMaker.Init()
	Group.GenerateID.Init()
}
