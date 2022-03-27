package assessment

import (
	uuid "github.com/satori/go.uuid"
	"github.com/yusufatalay/C1TakeHomeCase/database"
	"gorm.io/gorm"
)

type Assessment struct {
	AssessmentUUID string `gorm:"primarykey" json:"assessment_uuid"`
	QuestionCount  int    `gorm:"default:15" json:"question_count"`
	Duration       int    `gorm:"default:60" json:"duration"`
}

// init function automatically runs when this package imported
func init() {
	database.DBConn.AutoMigrate(&Assessment{})

}

func (a *Assessment) BeforeCreate(tx *gorm.DB) (err error) {

	a.AssessmentUUID = uuid.NewV4().String()
	return
}

// CreateAssessment creates an assessment with given parameters.
// returns UUID of the newly created assignment and an error if there is any.
func CreateAssessment() (string, error) {

	// for now assessments have hard coded (default) Duration and QuestionCount
	// variadic signature can be implemented later...
	var payload Assessment

	// save this assessment to db
	result := database.DBConn.Create(&payload)

	if result.Error != nil {
		return "", result.Error
	}

	return payload.AssessmentUUID, nil
}
