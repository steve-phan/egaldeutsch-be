package services

import (
	"egaldeutsch-be/modules/quiz/internal/models"
	"egaldeutsch-be/modules/quiz/internal/repositories"
)

type QuestionService struct {
	repo *repositories.QuestionRepository
}

func NewQuestionService(repo *repositories.QuestionRepository) *QuestionService {
	return &QuestionService{repo: repo}
}

func (s *QuestionService) CreateQuestion(question models.CreateQuestionDTO) error {

	// How do we pronounce this symbol?
	// How is this pronounced?
	q := &models.Question{
		QuestionText:  question.QuestionText,
		Options:       question.Options,
		CorrectOption: question.CorrectOption,
		Category:      question.Category,
	}
	return s.repo.Create(q)
}

func (s *QuestionService) GetAllQuestions() ([]models.Question, error) {
	return s.repo.GetAll()
}
