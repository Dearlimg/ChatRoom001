package dao

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
)

type database struct {
	DB    *sql.DB
	Redis *redis.Client
}

var Database = new(database)
