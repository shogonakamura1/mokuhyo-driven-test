package model

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty"`
}

type CreateProjectRequest struct {
	Title       string  `json:"title" binding:"required,min=3,max=20"`
	Description *string `json:"description,omitempty"`
}

type UpdateProjectRequest struct {
	Title       *string `json:"title,omitempty" binding:"omitempty,min=3,max=20"`
	Description *string `json:"description,omitempty"`
	Archived    *bool   `json:"archived,omitempty"`
}
