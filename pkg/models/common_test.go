package models

import (
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestBaseModelBeforeCreateSetsID(t *testing.T) {
	b := &BaseModel{}
	// Call BeforeCreate with nil tx (should not panic)
	if err := b.BeforeCreate(&gorm.DB{}); err != nil {
		t.Fatalf("BeforeCreate returned error: %v", err)
	}
	if b.ID == uuid.Nil {
		t.Fatalf("expected ID to be set")
	}
}
