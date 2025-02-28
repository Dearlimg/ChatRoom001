package dao

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
)

type database struct {
	DB    *sql.DB
	Redis *redis.Client
}

var Database = new(database)
