package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mokuhyo-driven-test/api/internal/model"
)

// UserRepository はユーザーリポジトリのインターフェースです
type UserRepository interface {
	GetByGoogleID(ctx context.Context, googleID string) (*model.User, error)
	Create(ctx context.Context, googleID, email, name string, picture *string) (*model.User, error)
	Update(ctx context.Context, userID uuid.UUID, email, name string, picture *string) (*model.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
}

// ProjectRepository はプロジェクトリポジトリのインターフェースです
type ProjectRepository interface {
	Create(ctx context.Context, userID uuid.UUID, req model.CreateProjectRequest) (*model.Project, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Project, error)
	GetByID(ctx context.Context, projectID uuid.UUID) (*model.Project, error)
	Update(ctx context.Context, projectID uuid.UUID, req model.UpdateProjectRequest) error
	CheckOwnership(ctx context.Context, projectID, userID uuid.UUID) (bool, error)
	UpdateUpdatedAt(ctx context.Context, projectID uuid.UUID) error
}

// NodeRepository はノードリポジトリのインターフェースです
type NodeRepository interface {
	Create(ctx context.Context, projectID uuid.UUID, content string, question *string) (*model.Node, error)
	GetByID(ctx context.Context, nodeID uuid.UUID) (*model.Node, error)
	Update(ctx context.Context, nodeID uuid.UUID, content string) error
	ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.Node, error)
	SoftDeleteWithDescendants(ctx context.Context, projectID, nodeID uuid.UUID) error
	GetMaxOrderIndex(ctx context.Context, projectID uuid.UUID, parentNodeID *uuid.UUID) (int, error)
}

// EdgeRepository はエッジリポジトリのインターフェースです
type EdgeRepository interface {
	Create(ctx context.Context, projectID uuid.UUID, parentNodeID *uuid.UUID, childNodeID uuid.UUID, relation model.RelationType, relationLabel *string, orderIndex int) (*model.Edge, error)
	ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.Edge, error)
	GetByID(ctx context.Context, edgeID uuid.UUID) (*model.Edge, error)
	Update(ctx context.Context, edgeID uuid.UUID, relation *string, relationLabel *string) error
	Reorder(ctx context.Context, projectID uuid.UUID, parentNodeID *uuid.UUID, orderedChildNodeIDs []uuid.UUID) error
}

// SettingsRepository は設定リポジトリのインターフェースです
type SettingsRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.UserSettings, error)
	Upsert(ctx context.Context, userID uuid.UUID, req model.UpdateSettingsRequest) (*model.UserSettings, error)
}
