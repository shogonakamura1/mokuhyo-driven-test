package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mokuhyo-driven-test/api/internal/model"
)

type EdgeRepository struct {
	db *DB
}

func NewEdgeRepository(db *DB) *EdgeRepository {
	return &EdgeRepository{db: db}
}

func (r *EdgeRepository) Create(ctx context.Context, projectID uuid.UUID, parentNodeID *uuid.UUID, childNodeID uuid.UUID, relation model.RelationType, relationLabel *string, orderIndex int) (*model.Edge, error) {
	var edge model.Edge
	err := r.db.pool.QueryRow(ctx, `
		INSERT INTO edges (project_id, parent_node_id, child_node_id, relation, relation_label, order_index)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, project_id, parent_node_id, child_node_id, relation, relation_label, order_index, created_at, updated_at
	`, projectID, parentNodeID, childNodeID, relation, relationLabel, orderIndex).Scan(
		&edge.ID, &edge.ProjectID, &edge.ParentNodeID, &edge.ChildNodeID,
		&edge.Relation, &edge.RelationLabel, &edge.OrderIndex,
		&edge.CreatedAt, &edge.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create edge: %w", err)
	}
	return &edge, nil
}

func (r *EdgeRepository) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.Edge, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT e.id, e.project_id, e.parent_node_id, e.child_node_id, e.relation, e.relation_label, e.order_index, e.created_at, e.updated_at
		FROM edges e
		INNER JOIN nodes n ON e.child_node_id = n.id
		WHERE e.project_id = $1 AND n.deleted_at IS NULL
		ORDER BY e.order_index
	`, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list edges: %w", err)
	}
	defer rows.Close()

	var edges []model.Edge
	for rows.Next() {
		var e model.Edge
		if err := rows.Scan(&e.ID, &e.ProjectID, &e.ParentNodeID, &e.ChildNodeID,
			&e.Relation, &e.RelationLabel, &e.OrderIndex,
			&e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan edge: %w", err)
		}
		edges = append(edges, e)
	}
	return edges, nil
}

func (r *EdgeRepository) GetByID(ctx context.Context, edgeID uuid.UUID) (*model.Edge, error) {
	var edge model.Edge
	err := r.db.pool.QueryRow(ctx, `
		SELECT id, project_id, parent_node_id, child_node_id, relation, relation_label, order_index, created_at, updated_at
		FROM edges
		WHERE id = $1
	`, edgeID).Scan(
		&edge.ID, &edge.ProjectID, &edge.ParentNodeID, &edge.ChildNodeID,
		&edge.Relation, &edge.RelationLabel, &edge.OrderIndex,
		&edge.CreatedAt, &edge.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get edge: %w", err)
	}
	return &edge, nil
}

func (r *EdgeRepository) Update(ctx context.Context, edgeID uuid.UUID, relation *string, relationLabel *string) error {
	query := "UPDATE edges SET updated_at = NOW()"
	args := []interface{}{}
	argIndex := 1

	if relation != nil {
		query += fmt.Sprintf(", relation = $%d", argIndex)
		args = append(args, *relation)
		argIndex++
	}
	if relationLabel != nil {
		query += fmt.Sprintf(", relation_label = $%d", argIndex)
		args = append(args, *relationLabel)
		argIndex++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argIndex)
	args = append(args, edgeID)

	_, err := r.db.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update edge: %w", err)
	}
	return nil
}

func (r *EdgeRepository) Reorder(ctx context.Context, projectID uuid.UUID, parentNodeID *uuid.UUID, orderedChildNodeIDs []uuid.UUID) error {
	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for i, childID := range orderedChildNodeIDs {
		query := `
			UPDATE edges
			SET order_index = $1, updated_at = NOW()
			WHERE project_id = $2 AND child_node_id = $3
		`
		args := []interface{}{i, projectID, childID}

		if parentNodeID == nil {
			query += " AND parent_node_id IS NULL"
		} else {
			query += " AND parent_node_id = $4"
			args = append(args, *parentNodeID)
		}

		_, err := tx.Exec(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to reorder edge: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
