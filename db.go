package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var conn *gorm.DB

// DBConnect starts up db connection
func DBConnect() *gorm.DB {
	fmt.Println("[DB] Connecting to db")

	if conn != nil {
		fmt.Println("[DB] Using existing connection")
		return conn
	}

	fmt.Println("[DB] Creating new connection")
	db, err := gorm.Open("sqlite3", "./db.db")
	if err != nil {
		panic("failed to connect database")
	}

	conn = db
	conn.AutoMigrate(&Project{})
	fmt.Println("[DB] Created, Migrated 🎉")

	return conn
}

// DBClose closes the database connection
func DBClose() {
	conn.Close()
	conn = nil
}
