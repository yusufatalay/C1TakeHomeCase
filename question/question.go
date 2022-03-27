package question

import (
	"github.com/yusufatalay/C1TakeHomeCase/category"
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	Text         string            `json:"text"`
	QuestionType string            `json:"question_type"`
	CategoryID   int               `json:"category_id"`
	Category     category.Category `json:"category"`
}
