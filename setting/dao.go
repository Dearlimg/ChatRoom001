package setting

import (
	"ChatRoom001/dao/mysql"
	"ChatRoom001/dao/redis"
)

type database struct {
}

func (d database) Init() {
	mysql.Init("root:123456@tcp(123.249.32.125:3306)/ChatRoom?charset=utf8mb4&parseTime=true")
	addr := "123.249.32.125:6381" // Redis 集群的地址
	password := "1234"            // Redis 密码
	poolSize := 10                // 连接池大小
	db := 0                       // 选择的数据库（默认是 0）
	redis.Init(addr, password, poolSize, db)

}
