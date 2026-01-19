package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mokuhyo-driven-test/api/internal/model"
	"github.com/mokuhyo-driven-test/api/internal/repository"
)

// projectRepository はプロジェクトリポジトリのPostgreSQL実装です
type projectRepository struct {
	db repository.DBInterface
}

// NewProjectRepository は新しいプロジェクトリポジトリを作成します
func NewProjectRepository(db repository.DBInterface) repository.ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, userID uuid.UUID, req model.CreateProjectRequest) (*model.Project, error) {
	var project model.Project
	err := r.db.QueryRow(ctx, `
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

func (r *projectRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Project, error) {
	rows, err := r.db.Query(ctx, `
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

func (r *projectRepository) GetByID(ctx context.Context, projectID uuid.UUID) (*model.Project, error) {
	var project model.Project
	err := r.db.QueryRow(ctx, `
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

func (r *projectRepository) Update(ctx context.Context, projectID uuid.UUID, req model.UpdateProjectRequest) error {
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

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}
	return nil
}

func (r *projectRepository) CheckOwnership(ctx context.Context, projectID, userID uuid.UUID) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM projects WHERE id = $1 AND user_id = $2
	`, projectID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check ownership: %w", err)
	}
	return count > 0, nil
}

func (r *projectRepository) UpdateUpdatedAt(ctx context.Context, projectID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE projects SET updated_at = NOW() WHERE id = $1
	`, projectID)
	if err != nil {
		return fmt.Errorf("failed to update updated_at: %w", err)
	}
	return nil
}

// nullString は *string を sql.NullString に変換します
func nullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// nullTime は *time.Time を sql.NullTime に変換します
func nullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

// stringPtr は sql.NullString を *string に変換します
func stringPtr(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// timePtr は sql.NullTime を *time.Time に変換します
func timePtr(t sql.NullTime) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
