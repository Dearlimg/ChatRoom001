package mysql

import (
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/dao/mysql/tx"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type DB interface {
	db.Querier
	tx.TXer
}

func Init(dataSourceName string) DB {
	// 建立连接
	pool, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect database: %v", err))
	}
	// 验证连接是否有效
	if err = pool.Ping(); err != nil {
		panic(fmt.Sprintf("Database connection verification failed: %v", err))
	}
	// 配置连接池（推荐生产环境参数）
	pool.SetMaxOpenConns(100)                // 最大打开连接数
	pool.SetMaxIdleConns(20)                 // 最大空闲连接数
	pool.SetConnMaxLifetime(5 * time.Minute) // 连接最大存活时间
	pool.SetConnMaxIdleTime(2 * time.Minute) // 空闲连接最大存活时间
	return &tx.SqlStore{Queries: db.New(pool), Pool: pool}
}
