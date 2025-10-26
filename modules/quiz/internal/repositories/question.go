package repositories

import (
	questionModels "egaldeutsch-be/modules/quiz/internal/models"

	"gorm.io/gorm"
)

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) Create(question *questionModels.Question) error {
	return r.db.Create(question).Error
}

func (r *QuestionRepository) GetAll() ([]questionModels.Question, error) {
	var questions []questionModels.Question
	err := r.db.Find(&questions).Error
	return questions, err
}
