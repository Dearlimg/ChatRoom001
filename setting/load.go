package setting

import (
	"ChatRoom001/dao"
	"ChatRoom001/global"
	"context"
)

type load struct {
}

//func(load )Init(){
//	var err error
//	err=tool.DoThat()
//}

func (load) LoadAllEmailToRedis() error {
	emails, err := dao.Database.DB.GetAllEmail(context.Background())
	if err != nil {
		return err
	}
	err = dao.Database.Redis.ReloadEmails(context.Background(), emails...)
	if err != nil {
		return err
	}
	global.Logger.Info("邮箱加载完成")
	return nil
}
