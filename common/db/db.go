package db

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var defaultdb *gorm.DB

// New initialze gorm db instance.
func New(dsn string) (*gorm.DB, error) {
	conn, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql db failed, %v", err)
	}

	conn.DB().Ping()
	conn.DB().SetConnMaxLifetime(time.Minute * 5)
	conn.DB().SetMaxIdleConns(10)
	conn.DB().SetMaxOpenConns(10)

	if defaultdb == nil {
		defaultdb = conn
	}

	return conn, nil
}

// Default returns default db instance.
func Default() *gorm.DB {
	return defaultdb
}
