package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mokuhyo-driven-test/api/internal/model"
	"github.com/mokuhyo-driven-test/api/internal/repository"
)

// userRepository はユーザーリポジトリのPostgreSQL実装です
type userRepository struct {
	db repository.DBInterface
}

// NewUserRepository は新しいユーザーリポジトリを作成します
func NewUserRepository(db repository.DBInterface) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, `
		SELECT id, google_id, email, name, picture, created_at, updated_at
		FROM users
		WHERE google_id = $1
	`, googleID).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.Picture,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by google ID: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, googleID, email, name string, picture *string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, `
		INSERT INTO users (google_id, email, name, picture)
		VALUES ($1, $2, $3, $4)
		RETURNING id, google_id, email, name, picture, created_at, updated_at
	`, googleID, email, name, picture).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.Picture,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, userID uuid.UUID, email, name string, picture *string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, `
		UPDATE users
		SET email = $1, name = $2, picture = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING id, google_id, email, name, picture, created_at, updated_at
	`, email, name, picture, userID).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.Picture,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, `
		SELECT id, google_id, email, name, picture, created_at, updated_at
		FROM users
		WHERE id = $1
	`, userID).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.Picture,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}
