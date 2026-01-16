package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mokuhyo-driven-test/api/internal/model"
)

type ProjectRepository struct {
	db *DB
}

func NewProjectRepository(db *DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, userID uuid.UUID, req model.CreateProjectRequest) (*model.Project, error) {
	var project model.Project
	err := r.db.pool.QueryRow(ctx, `
		INSERT INTO projects (user_id, title, description)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, title, description, created_at, updated_at, archived_at
	`, userID, req.Title, req.Description).Scan(
		&project.ID, &project.UserID, &project.Title, &project.Description,
		&project.CreatedAt, &project.UpdatedAt, &project.ArchivedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	return &project, nil
}

func (r *ProjectRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Project, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, user_id, title, description, created_at, updated_at, archived_at
		FROM projects
		WHERE user_id = $1 AND archived_at IS NULL
		ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Description,
			&p.CreatedAt, &p.UpdatedAt, &p.ArchivedAt); err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *ProjectRepository) GetByID(ctx context.Context, projectID uuid.UUID) (*model.Project, error) {
	var project model.Project
	err := r.db.pool.QueryRow(ctx, `
		SELECT id, user_id, title, description, created_at, updated_at, archived_at
		FROM projects
		WHERE id = $1
	`, projectID).Scan(
		&project.ID, &project.UserID, &project.Title, &project.Description,
		&project.CreatedAt, &project.UpdatedAt, &project.ArchivedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	return &project, nil
}

func (r *ProjectRepository) Update(ctx context.Context, projectID uuid.UUID, req model.UpdateProjectRequest) error {
	query := "UPDATE projects SET updated_at = NOW()"
	args := []interface{}{}
	argIndex := 1

	if req.Title != nil {
		query += fmt.Sprintf(", title = $%d", argIndex)
		args = append(args, *req.Title)
		argIndex++
	}
	if req.Description != nil {
		query += fmt.Sprintf(", description = $%d", argIndex)
		args = append(args, *req.Description)
		argIndex++
	}
	if req.Archived != nil {
		if *req.Archived {
			query += fmt.Sprintf(", archived_at = NOW()")
		} else {
			query += fmt.Sprintf(", archived_at = NULL")
		}
	}

	query += fmt.Sprintf(" WHERE id = $%d", argIndex)
	args = append(args, projectID)

	_, err := r.db.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}
	return nil
}

func (r *ProjectRepository) CheckOwnership(ctx context.Context, projectID, userID uuid.UUID) (bool, error) {
	var count int
	err := r.db.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM projects WHERE id = $1 AND user_id = $2
	`, projectID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check ownership: %w", err)
	}
	return count > 0, nil
}

func (r *ProjectRepository) UpdateUpdatedAt(ctx context.Context, projectID uuid.UUID) error {
	_, err := r.db.pool.Exec(ctx, `
		UPDATE projects SET updated_at = NOW() WHERE id = $1
	`, projectID)
	if err != nil {
		return fmt.Errorf("failed to update updated_at: %w", err)
	}
	return nil
}
