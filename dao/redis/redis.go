package redis

import (
	"ChatRoom001/dao/redis/operate"
	"ChatRoom001/model/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

//func Init(Addr []string, Password string, PoolSize int) *redis.ClusterClient {
//	rdb := redis.NewClusterClient(&redis.ClusterOptions{
//		Addrs:    Addr,     // 使用传入的地址
//		Password: Password, // 密码
//		PoolSize: PoolSize, // 连接池大小
//
//	})
//
//	// 测试连接
//	_, err := rdb.Ping(context.Background()).Result()
//	if err != nil {
//		log.Fatalf("Failed to connect to Redis cluster: %v", err)
//	}
//
//	return rdb
//}

func testredis() {
	// 集群模式配置
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"123.249.32.125:6381",
			"123.249.32.125:6382",
			"123.249.32.125:6383",
		},
		Password:     "1234",          // 集群密码
		PoolSize:     20,              // 连接池大小
		ReadOnly:     false,           // 是否使用只读节点
		DialTimeout:  5 * time.Second, // 连接超时
		ReadTimeout:  3 * time.Second, // 读取超时
		WriteTimeout: 3 * time.Second, // 写入超时
	})

	// 测试连接
	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Redis集群连接成功:", pong)
	rdb.SAdd(ctx, "chatroom", 1)
	rdb.SAdd(ctx, "EmailKey", 123)

}

func Init(cluster config.RedisCluster) *operate.RDB {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"123.249.32.125:6381",
			"123.249.32.125:6382",
			"123.249.32.125:6383",
		},
		Password:     "1234",          // 集群密码
		PoolSize:     20,              // 连接池大小
		ReadOnly:     false,           // 是否使用只读节点
		DialTimeout:  5 * time.Second, // 连接超时
		ReadTimeout:  3 * time.Second, // 读取超时
		WriteTimeout: 3 * time.Second, // 写入超时
	})

	// 测试连接
	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Redis集群连接成功:", pong)
	rdb.SAdd(ctx, "chatroom", 1)
	rdb.SAdd(ctx, "EmailKey", 123112)

	return operate.New(rdb)
}
