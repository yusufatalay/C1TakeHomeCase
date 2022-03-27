package models

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/yusufatalay/C1TakeHomeCase/database"
	"gorm.io/gorm"
)

// Answer has a 'belongs to' relationship with the Question
// Answer has a 'belongs to' relationship with the User
// Answer has a 'has one" relationship with the Option
type Answer struct {
	gorm.Model
	QuestionID uint     `json:"question_id"`
	Question   Question `json:"question"`
	UserID     uint     `json:"user_id"`
	User       User     `json:"user"`
	Text       string   `json:"text"`
	IsCorrect  bool     `json:"is_correct"`
	Option     Option   `json:"option"`
}

// SubmitAnswer submits user's answer to the given question and sets the other fields accordingly
func SubmitAnswer(c *fiber.Ctx) (err error) {

	// Unmarshall the body
	// I assume the body contains QuestionID, UserID, Option, Text (if the question is a short answer question), and other necessary fields
	var payload Answer
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	// update the option's isAnswer and answerid  fields
	var option Option
	database.DBConn.First(&option, payload.Option.ID)

	database.DBConn.Model(&option).Update("is_answer", 1)
	database.DBConn.Model(&option).Update("answer_id", payload.ID)

	// set the payload's iscorrect field
	payload.IsCorrect = option.IsCorrect

	// get the assessment and the question number for checking whether the user is at the last question or not
	assessmentuuid := c.Params("assessmentuuid")
	questionnumber, _ := c.ParamsInt("questionnumber") // unhandled exception, I assume that the user hasn't modified the URL when submitting the answer

	var assessment Assessment
	database.DBConn.First(&assessment, "assessment_uuid = ?", assessmentuuid)

	if questionnumber == assessment.QuestionCount {

		c.JSON(Response{
			Data:    payload,
			Message: "Last question has submitted",
			Type:    true,
		})
		return
	}
	// write the answer to the database
	database.DBConn.Create(&payload)

	c.JSON(Response{
		Data:    payload,
		Message: fmt.Sprintf("%d. Question answered", questionnumber),
		Type:    true,
	})
	return
}
