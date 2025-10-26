package hanlers

import (
	"egaldeutsch-be/modules/quiz/internal/models"
	"egaldeutsch-be/modules/quiz/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	service *services.QuestionService
}

func NewQuestionHandler(service *services.QuestionService) *QuestionHandler {
	return &QuestionHandler{service: service}
}

func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var req models.CreateQuestionDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateQuestion(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Question created successfully"})
}

func (h *QuestionHandler) GetAllQuestions(c *gin.Context) {
	questions, err := h.service.GetAllQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}
