package redis

import (
	"context"
	"log"
	"github.com/go-redis/redis/v8"
)

func Init(Addr, Password string, PoolSize, DB int) *redis.ClusterClient {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{Addr}, // Redis 集群的节点地址列表
		Password: Password,       // 密码
		PoolSize: PoolSize,       // 连接池大小
	})

	// 测试连接
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis cluster: %v", err)
	}

	return rdb
}
