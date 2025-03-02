package dao

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
)

type database struct {
	DB    *sql.DB
	Redis *redis.ClusterClient
}

var Database = new(database)
