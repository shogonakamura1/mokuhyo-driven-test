package supabase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mokuhyo-driven-test/api/internal/model"
	"github.com/mokuhyo-driven-test/api/internal/repository"
)

// settingsRepository は設定リポジトリのSupabase実装です
type settingsRepository struct {
	db repository.DBInterface
}

// NewSettingsRepository は新しい設定リポジトリを作成します
func NewSettingsRepository(db repository.DBInterface) repository.SettingsRepository {
	return &settingsRepository{db: db}
}

func (r *settingsRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.UserSettings, error) {
	var settings model.UserSettings
	err := r.db.QueryRow(ctx, `
		SELECT user_id, theme, accent_color, updated_at
		FROM user_settings
		WHERE user_id = $1
	`, userID).Scan(
		&settings.UserID, &settings.Theme, &settings.AccentColor, &settings.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		// Return default settings
		return &model.UserSettings{
			UserID:      userID,
			Theme:       "light",
			AccentColor: "blue",
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}
	return &settings, nil
}

func (r *settingsRepository) Upsert(ctx context.Context, userID uuid.UUID, req model.UpdateSettingsRequest) (*model.UserSettings, error) {
	// Get current settings or create default
	current, _ := r.GetByUserID(ctx, userID)

	theme := current.Theme
	if req.Theme != nil {
		theme = *req.Theme
	}

	accentColor := current.AccentColor
	if req.AccentColor != nil {
		accentColor = *req.AccentColor
	}

	_, err := r.db.Exec(ctx, `
		INSERT INTO user_settings (user_id, theme, accent_color, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id) DO UPDATE
		SET theme = $2, accent_color = $3, updated_at = NOW()
	`, userID, theme, accentColor)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert settings: %w", err)
	}

	return &model.UserSettings{
		UserID:      userID,
		Theme:       theme,
		AccentColor: accentColor,
	}, nil
}
