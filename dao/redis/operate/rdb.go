package operate

import (
	"github.com/redis/go-redis/v9"
)

type RDB struct {
	rdb *redis.ClusterClient
}

func New(rdb *redis.ClusterClient) *RDB {
	return &RDB{rdb: rdb}
}
