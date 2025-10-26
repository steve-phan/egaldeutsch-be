package quiz

import (
	"egaldeutsch-be/modules/quiz/internal/hanlers"
	"egaldeutsch-be/modules/quiz/internal/repositories"
	"egaldeutsch-be/modules/quiz/internal/services"

	"gorm.io/gorm"
)

type Module struct {
	handler *hanlers.QuestionHandler
	service *services.QuestionService
	repo    *repositories.QuestionRepository
}

func NewModule(db *gorm.DB) *Module {

	repo := repositories.NewQuestionRepository(db)

	service := services.NewQuestionService(repo)

	handler := hanlers.NewQuestionHandler(service)

	return &Module{handler: handler, service: service, repo: repo}
}
