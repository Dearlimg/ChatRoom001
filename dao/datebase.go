package dao

import (
	"ChatRoom001/dao/mysql"
	"ChatRoom001/dao/redis/operate"
)

type database struct {
	DB    mysql.DB
	Redis *operate.RDB
}

var Database = new(database)
