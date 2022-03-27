package models

import (
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yusufatalay/C1TakeHomeCase/database"
)

type User struct {
	ID        uint `gorm:"primarykey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time `json:"-"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email"`
}

// CreateUser creates the user with given information (FirstName, LastName, email) this informations are gathered from icoming request's body (POST)
func CreateUser(c *fiber.Ctx) error {
	// unmarshall the reques's body into the payload (User object (kind of...))
	var payload User
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	// write user to the database
	database.DBConn.Create(&payload)

	err := c.Status(http.StatusCreated).JSON(payload)

	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(Response{
		Data:    payload,
		Message: "User has created and inserted to database",
		Type:    true,
	})

}
