package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID   uuid.UUID `gorm:"column:id"`
	Name string    `gorm:"column:name"`
	Role string    `gorm:"column:role"` // "admin" or "user"
}

// QuizQuestion represents a single question in a quiz
type QuizQuestion struct {
	ID      uuid.UUID `json:"id" db:"id"`
	Prompt  string    `json:"prompt" db:"prompt"`
	Options []string  `json:"options" db:"options"`
	Answer  int       `json:"answer" db:"answer"`
}

// Quiz represents a quiz with multiple questions
type Quiz struct {
	ID        uuid.UUID      `json:"id" db:"id"`
	Title     string         `json:"title" db:"title"`
	Questions []QuizQuestion `json:"questions" db:"questions"`
}

// Article represents a learning article
type Article struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Summary   string    `json:"summary" db:"summary"`
	Content   string    `json:"content" db:"content"`
	Level     string    `json:"level" db:"level"` // "A1", "A2", "B1", "B2", "C1"
	AuthorID  uuid.UUID `json:"author_id" db:"author_id"`
	Author    User      `json:"author" db:"author"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Quiz      Quiz      `json:"quiz,omitempty" db:"quiz"`
}

// ArticleData represents the data needed to create/update an article
type ArticleData struct {
	Title    string    `json:"title" binding:"required"`
	Summary  string    `json:"summary" binding:"required"`
	Content  string    `json:"content" binding:"required"`
	Level    string    `json:"level" binding:"required,oneof=A1 A2 B1 B2 C1"`
	AuthorID uuid.UUID `json:"author_id" binding:"required"`
	Quiz     *Quiz     `json:"quiz,omitempty"`
}

// PaginatedArticles represents a paginated response of articles
type PaginatedArticles struct {
	Items      []Article `json:"items"`
	Page       int       `json:"page"`
	PerPage    int       `json:"per_page"`
	TotalItems int64     `json:"total_items"`
	TotalPages int       `json:"total_pages"`
}

// CreateArticleRequest represents the request payload for creating an article
type CreateArticleRequest struct {
	Article ArticleData `json:"article" binding:"required"`
}

// UpdateArticleRequest represents the request payload for updating an article
type UpdateArticleRequest struct {
	Article ArticleData `json:"article" binding:"required"`
}
