package setting

type group struct {
	Config config
	Dao    database
	Logger log
}

var Group group

func Inits() {
	Group.Config.Init()
	Group.Dao.Init()
	Group.Logger.Init()
}
