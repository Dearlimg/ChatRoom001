package mysql

import (
	db "ChatRoom001/dao/mysql/sqlc"
	"ChatRoom001/dao/mysql/tx"
)

//var db *sql.DB

type DB interface {
	db.Querier
	tx.TXer
}

//func Init(dataSourceName string) DB {
//	// Open a connection to the database
//	db, err := sql.Open("mysql", dataSourceName)
//	if err != nil {
//		panic(err)
//	}
//	// Check if the connection is alive
//	err = db.Ping()
//	db.SetMaxOpenConns(10)           // Set max open connections
//	db.SetMaxIdleConns(5)            // Set max idle connections
//	db.SetConnMaxLifetime(time.Hour) // Set max connection lifetime
//
//	return &tx.Sql
//}
