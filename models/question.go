package models

import (
	"gorm.io/gorm"
)

// Question has a "Belongs To" relationship with the Category
// Question has a "Belongs To" relationship with the Assessment
// Question has a "Has Many" relationship with the Options
type Question struct {
	gorm.Model
	Text          string `json:"text"`
	QuestionType  string `json:"question_type"`
	QuestionFile  string `json:"question_file"` // contains the name ("path") of the question file
	AnswerText    string `json:"answer_text"`   // if the question is a short-answer type of question (like 11th question in given example) this variable will hold the answer
	Point         uint   `json:"point"`
	RemainingTime int    `json:"remaining_time"` // in seconds, every question send to the front end will have this updated, will be calculated
	// relationship are setted here
	AssessmentUUID string   `json:"assessment_uuid"` // this is foreign key to the assessment that the quesiton belongs to
	Options        []Option `json:"options"`
	//CategoryID     int        `json:"category_id"`
	//Category       Category   `json:"category"`
}

// AfterUpdate assigns related options to the created question
// Questions will get updated whenever they get assigned to an assessment
func (q *Question) AfterUpdate(tx *gorm.DB) (err error) {
	// get options that has the created question's id
	var options []Option
	tx.Where("question_id = ?", q.ID).Find(&options)
	// append the options to the question
	q.Options = append(q.Options, options...)

	return tx.Error
}
