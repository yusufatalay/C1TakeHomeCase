package models

import (
	"time"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/yusufatalay/C1TakeHomeCase/database"
	"gorm.io/gorm"
)

// Assessment has a "Belongs To" relationship with User
// Assessment has a "Has Many" relationship with Question
type Assessment struct {
	AssessmentUUID string     `gorm:"primarykey" json:"assessment_uuid"`
	CreatedAt      time.Time  `json:"created_at"`
	QuestionCount  int        `gorm:"default:15" json:"question_count"`
	Duration       int        `gorm:"default:60" json:"duration"` // in minutes
	Questions      []Question `gorm:"foreignKey:AssessmentUUID" json:"questions"`
	UserID         uint       `json:"taker_id"`
	User           User       `json:"taker"`
}

// BeforeCreate is a trigger which adds an UUID to the assessment before inserting that into the database
func (a *Assessment) BeforeCreate(tx *gorm.DB) (err error) {

	a.AssessmentUUID = uuid.NewV4().String()
	return
}

// CreateAssessment creates an assessment with given parameters.
// returns UUID of the newly created assignment and an error if there is any.
func CreateAssessment(c *fiber.Ctx) (err error) {

	var payload Assessment
	var taker User
	// get desired users id form URL parameter
	takerid, err := c.ParamsInt("takerid")
	if err != nil {
		c.JSON(Response{
			Data:    nil,
			Message: "User ID must be a positive integer",
			Type:    false,
		})
	}
	database.DBConn.First(&taker, takerid)
	if err = c.BodyParser(&payload); err != nil {
		c.JSON(Response{
			Data:    nil,
			Message: "Database error",
			Type:    false,
		})
		return
	}
	// set the other fields
	payload.UserID = uint(takerid)
	payload.User = taker

	// save this assessment to db
	result := database.DBConn.Create(&payload)

	if result.Error != nil {
		c.JSON(Response{
			Data:    nil,
			Message: "Cannot write to the database",
			Type:    false,
		})
		return result.Error
	}
	// everything has gone well notify the user
	c.JSON(Response{
		Data:    payload,
		Message: "Assignment created",
		Type:    true,
	})
	return
}

func (a *Assessment) AfterCreate(tx *gorm.DB) (err error) {
	// make request to the database and get question amont of questions from it
	questions := make([]Question, a.QuestionCount)
	result := tx.Limit(a.QuestionCount).Find(&questions)
	for _, question := range questions {
		tx.Model(&question).Update("assessment_uuid", a.AssessmentUUID)
	}
	a.Questions = append(a.Questions, questions...)

	return result.Error
}

// GetAssessmentQuestion will return the question with matching id and assessment uuid -> url format GET /api/v1/<assessment's uuid>/<question's number (1- assessment.QuestionAmount)>
func GetAssessmentQuestion(c *fiber.Ctx) (err error) {
	// get assessment's uuid from the url
	assessmentuuid := c.Params("assessmentuuid")
	// get desired question's number
	questionnumber, err := c.ParamsInt("questionnumber")

	if err != nil {
		c.JSON(Response{
			Data:    nil,
			Message: "Question number must be a positive integer",
			Type:    false,
		})
		return
	}

	// check if that assessment exists
	var assessment Assessment
	result := database.DBConn.First(&assessment, "assessment_uuid = ?", assessmentuuid)
	// if assessment is not in the database , notify the user
	if result.Error != nil {
		c.JSON(Response{
			Data:    nil,
			Message: "Assessment cannot found",
			Type:    false,
		})
		return
	}

	// check if question number falls into the assessment's question amount range
	// NOTE: User should not be able to see the questionnumer'th question if the previous question hasn't been answered
	if questionnumber-1 > assessment.QuestionCount || questionnumber < 0 {
		c.JSON(Response{
			Data:    nil,
			Message: "Question number exceeds the amount of questions this assessment have",
			Type:    false,
		})
		return
	}
	// retrieve the questions with that assessment uuid
	var questions []Question
	result = database.DBConn.Where("assessment_uuid = ?", assessmentuuid).Find(&questions)

	if result.Error != nil {
		c.JSON(Response{
			Data:    nil,
			Message: "No question found with given assessment uuid",
			Type:    false,
		})
		return
	}
	// get question's options
	var options []Option
	result = database.DBConn.Where("question_id = ?", questions[questionnumber-1].ID).Find(&options)
	questions[questionnumber-1].Options = append(questions[questionnumber-1].Options, options...)
	if result.Error != nil {
		c.JSON(Response{
			Data:    nil,
			Message: "No option(s) found with given question id",
			Type:    false,
		})
		return
	}

	// set remaining time attribute of the question
	// Note: this calculation gives a very big number which is I guess not correct.
	questions[questionnumber-1].RemainingTime = (assessment.Duration)*60 - int(assessment.CreatedAt.Sub(time.Now()).Seconds())

	if questions[questionnumber-1].RemainingTime <= 0 {
		c.JSON(Response{
			Data:    nil,
			Message: "Time is up",
			Type:    true,
		})
	}

	// everything goes well then provide the question
	c.JSON(Response{
		Data:    questions[questionnumber-1],
		Message: "Question retrieved successfully",
		Type:    true,
	})

	return
}
