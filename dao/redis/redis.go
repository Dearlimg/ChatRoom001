package redis

import (
	"ChatRoom001/dao/redis/operate"
	"ChatRoom001/model/config"
	"github.com/redis/go-redis/v9"
)

// Init RedisCluster init 造成了一些参数浪费,ClusterOptions不支持,就当为日后扩展做准备吧
func Init(cluster config.RedisCluster) *operate.RDB {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        cluster.Endpoints,
		Password:     cluster.Password,
		PoolSize:     cluster.PoolConfig.PoolSize,
		ReadTimeout:  cluster.PoolConfig.ReadTimeout,
		WriteTimeout: cluster.PoolConfig.WriteTimeout,
		ReadOnly:     cluster.ReadOnly,
		DialTimeout:  cluster.PoolConfig.ConnectTimeout,
	})
	return operate.New(rdb)
}
