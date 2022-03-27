package models

import "gorm.io/gorm"

// Option has a "Belongs To" relationship with the Question
type Option struct {
	gorm.Model
	Text       string `json:"text"`
	IsAnswer   bool   `json:"is_answer"`
	IsCorrect  bool   `json:"-"`           // user should not see whether if this option is correct or not
	OptionFile string `json:"option_file"` // contains the name or path of the option file
	QuestionID uint   `json:"question_id"` // foreign key to the question
	AnswerID   uint   `json:"answer_id"`   // foreign key to the answer
}
