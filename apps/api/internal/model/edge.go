package model

import (
	"time"

	"github.com/google/uuid"
)

type RelationType string

const (
	RelationNeutral  RelationType = "neutral"
	RelationWhy      RelationType = "why"
	RelationConcrete RelationType = "concrete"
	RelationHow      RelationType = "how"
	RelationWhat     RelationType = "what"
	RelationCustom   RelationType = "custom"
)

type Edge struct {
	ID           uuid.UUID   `json:"id"`
	ProjectID    uuid.UUID   `json:"project_id"`
	ParentNodeID *uuid.UUID  `json:"parent_node_id"`
	ChildNodeID  uuid.UUID   `json:"child_node_id"`
	Relation     RelationType `json:"relation"`
	RelationLabel *string     `json:"relation_label,omitempty"`
	OrderIndex   int          `json:"order_index"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type UpdateEdgeRequest struct {
	Relation      *string `json:"relation,omitempty"`
	RelationLabel *string `json:"relation_label,omitempty" binding:"omitempty,max=20"`
}

type ReorderRequest struct {
	ParentNodeID        *uuid.UUID `json:"parent_node_id"`
	OrderedChildNodeIDs []uuid.UUID `json:"ordered_child_node_ids" binding:"required"`
}
