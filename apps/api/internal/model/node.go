package model

import (
	"time"

	"github.com/google/uuid"
)

type Node struct {
	ID        uuid.UUID  `json:"id"`
	ProjectID uuid.UUID  `json:"project_id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type CreateNodeRequest struct {
	Content       string     `json:"content" binding:"required,min=1,max=200"`
	ParentNodeID  *uuid.UUID `json:"parent_node_id"`
	Relation      string     `json:"relation,omitempty"`
	RelationLabel *string    `json:"relation_label,omitempty"`
	OrderIndex    *int       `json:"order_index,omitempty"`
}

type UpdateNodeRequest struct {
	Content string `json:"content" binding:"required,min=1,max=200"`
}
