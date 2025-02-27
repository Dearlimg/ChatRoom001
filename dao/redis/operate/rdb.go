package operate

import "github.com/go-redis/redis/v8"

type RDB struct {
	rdb *redis.ClusterClient
}

func New(rdb *redis.ClusterClient) *RDB {
	return &RDB{rdb: rdb}
}
