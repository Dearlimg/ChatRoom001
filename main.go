package main

import (
	"ChatRoom001/dao/mysql"
	"ChatRoom001/dao/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	mysql.Init("root:123456@tcp(123.249.32.125:3306)/ChatRoom?charset=utf8mb4&parseTime=true")
	Addrs := []string{
		"192.168.12.102:6381", // 正确的 Redis 地址
		"123.249.32.125:6382",
		"123.249.32.125:6383",
	}
	password := "1234" // Redis 密码
	poolSize := 10     // 连接池大小
	db := 0            // 选择的数据库（默认是 0）
	//redis.Init(Addrs, password, poolSize, db)
	redis.Init1(Addrs, password, poolSize, db)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 200,
		})
	})
	r.Run(":8080")

}
