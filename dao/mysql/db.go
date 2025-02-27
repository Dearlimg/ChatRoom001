package mysql

import (
	"database/sql"
	"fmt"
	"time"
)

var db *sql.DB

func Init(dsn string) (*sql.DB, error) {
	// Open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Check if the connection is alive
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	db.SetMaxOpenConns(10)           // Set max open connections
	db.SetMaxIdleConns(5)            // Set max idle connections
	db.SetConnMaxLifetime(time.Hour) // Set max connection lifetime

	return db, nil
}
