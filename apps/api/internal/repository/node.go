package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mokuhyo-driven-test/api/internal/model"
)

type NodeRepository struct {
	db *DB
}

func NewNodeRepository(db *DB) *NodeRepository {
	return &NodeRepository{db: db}
}

func (r *NodeRepository) Create(ctx context.Context, projectID uuid.UUID, content string) (*model.Node, error) {
	var node model.Node
	err := r.db.pool.QueryRow(ctx, `
		INSERT INTO nodes (project_id, content)
		VALUES ($1, $2)
		RETURNING id, project_id, content, created_at, updated_at, deleted_at
	`, projectID, content).Scan(
		&node.ID, &node.ProjectID, &node.Content,
		&node.CreatedAt, &node.UpdatedAt, &node.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create node: %w", err)
	}
	return &node, nil
}

func (r *NodeRepository) GetByID(ctx context.Context, nodeID uuid.UUID) (*model.Node, error) {
	var node model.Node
	err := r.db.pool.QueryRow(ctx, `
		SELECT id, project_id, content, created_at, updated_at, deleted_at
		FROM nodes
		WHERE id = $1 AND deleted_at IS NULL
	`, nodeID).Scan(
		&node.ID, &node.ProjectID, &node.Content,
		&node.CreatedAt, &node.UpdatedAt, &node.DeletedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}
	return &node, nil
}

func (r *NodeRepository) Update(ctx context.Context, nodeID uuid.UUID, content string) error {
	_, err := r.db.pool.Exec(ctx, `
		UPDATE nodes SET content = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`, content, nodeID)
	if err != nil {
		return fmt.Errorf("failed to update node: %w", err)
	}
	return nil
}

func (r *NodeRepository) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.Node, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, project_id, content, created_at, updated_at, deleted_at
		FROM nodes
		WHERE project_id = $1 AND deleted_at IS NULL
	`, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}
	defer rows.Close()

	var nodes []model.Node
	for rows.Next() {
		var n model.Node
		if err := rows.Scan(&n.ID, &n.ProjectID, &n.Content,
			&n.CreatedAt, &n.UpdatedAt, &n.DeletedAt); err != nil {
			return nil, fmt.Errorf("failed to scan node: %w", err)
		}
		nodes = append(nodes, n)
	}
	return nodes, nil
}

func (r *NodeRepository) SoftDeleteWithDescendants(ctx context.Context, projectID, nodeID uuid.UUID) error {
	_, err := r.db.pool.Exec(ctx, `
		WITH RECURSIVE subtree AS (
			SELECT e.child_node_id AS id
			FROM edges e
			WHERE e.project_id = $1 AND e.child_node_id = $2
			UNION ALL
			SELECT e.child_node_id
			FROM edges e
			JOIN subtree s ON e.parent_node_id = s.id
			WHERE e.project_id = $1
		)
		UPDATE nodes
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE project_id = $1
		  AND id IN (SELECT id FROM subtree)
	`, projectID, nodeID)
	if err != nil {
		return fmt.Errorf("failed to soft delete with descendants: %w", err)
	}
	return nil
}

func (r *NodeRepository) GetMaxOrderIndex(ctx context.Context, projectID uuid.UUID, parentNodeID *uuid.UUID) (int, error) {
	var maxOrder int
	query := `
		SELECT COALESCE(MAX(order_index), -1) + 1
		FROM edges
		WHERE project_id = $1
	`
	args := []interface{}{projectID}

	if parentNodeID == nil {
		query += " AND parent_node_id IS NULL"
	} else {
		query += " AND parent_node_id = $2"
		args = append(args, *parentNodeID)
	}

	err := r.db.pool.QueryRow(ctx, query, args...).Scan(&maxOrder)
	if err != nil {
		return 0, fmt.Errorf("failed to get max order index: %w", err)
	}
	return maxOrder, nil
}
