package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"egaldeutsch-be/internal/database"
	"egaldeutsch-be/internal/models"
)

type ArticleHandler struct {
	db *database.Database
}

func NewArticleHandler(db *database.Database) *ArticleHandler {
	return &ArticleHandler{db: db}
}

// GetArticles handles GET /api/v1/articles
func (h *ArticleHandler) GetArticles(c *gin.Context) {
	page := 1
	perPage := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if pp := c.Query("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 && parsed <= 100 {
			perPage = parsed
		}
	}

	articles, totalItems, err := h.getArticlesFromDB(page, perPage)
	if err != nil {
		logrus.WithError(err).Error("Failed to get articles")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve articles"})
		return
	}

	totalPages := int((totalItems + int64(perPage) - 1) / int64(perPage))

	response := models.PaginatedArticles{
		Items:      articles,
		Page:       page,
		PerPage:    perPage,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// GetArticle handles GET /api/v1/articles/:id
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	article, err := h.getArticleFromDB(id)
	if err != nil {
		logrus.WithError(err).WithField("article_id", id).Error("Failed to get article")
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, article)
}

// CreateArticle handles POST /api/v1/articles
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req models.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article := models.Article{
		ID:        uuid.New(),
		Title:     req.Article.Title,
		Summary:   req.Article.Summary,
		Content:   req.Article.Content,
		Level:     req.Article.Level,
		AuthorID:  req.Article.AuthorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if req.Article.Quiz != nil {
		article.Quiz = *req.Article.Quiz
	}

	if err := h.createArticleInDB(&article); err != nil {
		logrus.WithError(err).Error("Failed to create article")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	c.JSON(http.StatusCreated, article)
}

// UpdateArticle handles PUT /api/v1/articles/:id
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var req models.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article := models.Article{
		ID:        id,
		Title:     req.Article.Title,
		Summary:   req.Article.Summary,
		Content:   req.Article.Content,
		Level:     req.Article.Level,
		AuthorID:  req.Article.AuthorID,
		UpdatedAt: time.Now(),
	}

	if req.Article.Quiz != nil {
		article.Quiz = *req.Article.Quiz
	}

	if err := h.updateArticleInDB(&article); err != nil {
		logrus.WithError(err).WithField("article_id", id).Error("Failed to update article")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	c.JSON(http.StatusOK, article)
}

// DeleteArticle handles DELETE /api/v1/articles/:id
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	if err := h.deleteArticleFromDB(id); err != nil {
		logrus.WithError(err).WithField("article_id", id).Error("Failed to delete article")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

// Database operations (these would be implemented with actual SQL queries)
func (h *ArticleHandler) getArticlesFromDB(page, perPage int) ([]models.Article, int64, error) {
	// TODO: Implement actual database query
	// This is a placeholder that returns mock data
	return []models.Article{}, 0, nil
}

func (h *ArticleHandler) getArticleFromDB(id uuid.UUID) (*models.Article, error) {
	// TODO: Implement actual database query
	// This is a placeholder
	return nil, fmt.Errorf("article not found")
}

func (h *ArticleHandler) createArticleInDB(article *models.Article) error {
	// TODO: Implement actual database insert
	// This is a placeholder
	return nil
}

func (h *ArticleHandler) updateArticleInDB(article *models.Article) error {
	// TODO: Implement actual database update
	// This is a placeholder
	return nil
}

func (h *ArticleHandler) deleteArticleFromDB(id uuid.UUID) error {
	// TODO: Implement actual database delete
	// This is a placeholder
	return nil
}
