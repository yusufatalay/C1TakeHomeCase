package userassessment

import (
	"github.com/yusufatalay/C1TakeHomeCase/database"
	"gorm.io/gorm"
)

type UserAssessment struct {
	gorm.Model
	AssessmentUUID string `json:"assessment"`
	UserID         int    `json:"exam_taker"`
}

func init() {
	database.DBConn.AutoMigrate(&UserAssessment{})
}

// CreateUserAssessment connects User and the Assessments table
func CreateUserAssessment(userid int, assessmentuuid string) error {
	var payload UserAssessment

	payload.AssessmentUUID = assessmentuuid
	payload.UserID = userid

	result := database.DBConn.Create(&payload)

	return result.Error
}
