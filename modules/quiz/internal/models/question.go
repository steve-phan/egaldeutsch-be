package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"egaldeutsch-be/pkg/models"
)

// Options is a custom type for a string slice to handle JSONB serialization.
type Options []string

// Value implements the driver.Valuer interface. This method is called when
// preparing the value for the database. It marshals the string slice into a
// JSON byte slice.
func (o Options) Value() (driver.Value, error) {
	if len(o) == 0 {
		// Return an empty JSON array if the slice is empty
		return json.Marshal([]string{})
	}
	return json.Marshal(o)
}

// Scan implements the sql.Scanner interface. This method is called when
// reading a value from the database. It unmarshals a JSON byte slice into
// the string slice.
func (o *Options) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	return json.Unmarshal(bytes, o)
}

type Question struct {
	models.BaseModel
	QuestionText  string  `json:"question_text" gorm:"not null;size:500"`
	Options       Options `json:"options" gorm:"type:jsonb;not null"`
	CorrectOption int     `json:"correct_option" gorm:"not null"`
	Category      string  `json:"category" gorm:"not null;size:100"`
}
