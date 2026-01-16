package model

import (
	"time"

	"github.com/google/uuid"
)

type UserSettings struct {
	UserID      uuid.UUID `json:"user_id"`
	Theme       string    `json:"theme"`
	AccentColor string    `json:"accent_color"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateSettingsRequest struct {
	Theme       *string `json:"theme,omitempty" binding:"omitempty,oneof=light dark"`
	AccentColor *string `json:"accent_color,omitempty"`
}

type MeResponse struct {
	User     UserInfo     `json:"user"`
	Settings UserSettings `json:"settings"`
}

type UserInfo struct {
	ID uuid.UUID `json:"id"`
}
