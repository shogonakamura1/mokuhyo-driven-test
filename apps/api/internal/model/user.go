package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	GoogleID    string    `json:"google_id" db:"google_id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Picture     *string   `json:"picture,omitempty"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
