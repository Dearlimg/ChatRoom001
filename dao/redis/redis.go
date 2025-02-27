package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

func Init1(Addr []string, Password string, PoolSize, DB int) *redis.ClusterClient {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    Addr,     // 使用传入的地址
		Password: Password, // 密码
		PoolSize: PoolSize, // 连接池大小
	})

	// 测试连接
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis cluster: %v", err)
	}

	return rdb
}
