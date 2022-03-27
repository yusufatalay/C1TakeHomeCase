package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
)

func init() {
	var err error
	DBConn, err = gorm.Open(sqlite.Open("databse.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}
	log.Println("Database connection has been made")

	// Do database migrations
	log.Println("Database migrations completed")
}
