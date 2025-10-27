package models

type CreateQuestionDTO struct {
	QuestionText  string   `json:"question_text" binding:"required,max=500"`
	Options       []string `json:"options" binding:"required,min=2,dive,required"`
	CorrectOption int      `json:"correct_option" binding:"gte=0"`
	Category      string   `json:"category" binding:"required,max=100"`
}
