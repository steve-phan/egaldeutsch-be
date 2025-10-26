package models

import "egaldeutsch-be/pkg/models"

type Question struct {
	models.BaseModel
	QuestionText  string   `json:"question_text" gorm:"not null;size:500"`
	Options       []string `json:"options" gorm:"type:jsonb;not null"`
	CorrectOption int      `json:"correct_option" gorm:"not null"`
	Category      string   `json:"category" gorm:"not null;size:100"`
}
