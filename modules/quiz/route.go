package quiz

import "github.com/gin-gonic/gin"

func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	quizs := r.Group("/quiz")

	{
		quizs.POST("/questions", m.handler.CreateQuestion)
		quizs.GET("/questions", m.handler.GetAllQuestions)
	}
}
