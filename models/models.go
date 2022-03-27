package models

import (
	"log"

	"github.com/yusufatalay/C1TakeHomeCase/database"
)

// this particular file will be responsible for Automigration of the other models

func init() {
	database.DBConn.AutoMigrate(&User{}, &Assessment{}, &Question{}, &Option{}, &Answer{}, &Feedback{})
	log.Println("Database migrations are completed")
}

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Type    bool        `json:"type"`
}
