package setting

import (
	"ChatRoom001/dao"
	"ChatRoom001/dao/mysql"
	"ChatRoom001/dao/redis"
	"ChatRoom001/global"
)

type database struct {
}

func (d database) Init() {
	dao.Database.DB = mysql.Init(global.PrivateSetting.Mysql.SourceName)
	dao.Database.Redis = redis.Init(global.PrivateSetting.RedisCluster)
}
