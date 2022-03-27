package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yusufatalay/C1TakeHomeCase/database"
	"gorm.io/gorm"
)

// Feedback has a "belongs to" relationship with User
type Feedback struct {
	gorm.Model
	Feedback string `json:"feedback"`
	UserID   uint   `json:"user_id"`
	User     User   `json:"user"`
}

func SubmitFeedback(c *fiber.Ctx) (err error) {
	var payload Feedback
	err = c.BodyParser(&payload)

	if err != nil {
		c.JSON(Response{
			Data:    nil,
			Message: "Could not submit feedback",
			Type:    false,
		})
		return
	}
	var user User
	database.DBConn.First(&user, payload.UserID)
	payload.User = user

	result := database.DBConn.Create(&payload)

	if result.Error != nil {
		c.JSON(Response{
			Data:    nil,
			Message: "Could not submit the feedback",
			Type:    false,
		})
		return result.Error
	}

	c.JSON(Response{
		Data:    payload,
		Message: "Thanks for your feedback",
		Type:    true,
	})
	return
}
