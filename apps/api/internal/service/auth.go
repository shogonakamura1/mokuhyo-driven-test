package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mokuhyo-driven-test/api/internal/model"
	"github.com/mokuhyo-driven-test/api/internal/repository"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// GetOrCreateUser はGoogle IDでユーザーを取得し、存在しない場合は作成します
func (s *AuthService) GetOrCreateUser(ctx context.Context, googleID, email, name string, picture *string) (*model.User, error) {
	// 既存のユーザーを取得
	user, err := s.userRepo.GetByGoogleID(ctx, googleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// ユーザーが存在する場合は更新して返す
	if user != nil {
		// 情報が変更されている可能性があるので更新
		updatedUser, err := s.userRepo.Update(ctx, user.ID, email, name, picture)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
		return updatedUser, nil
	}

	// ユーザーが存在しない場合は作成
	newUser, err := s.userRepo.Create(ctx, googleID, email, name, picture)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

// GetUserByID はIDでユーザーを取得します
func (s *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByGoogleID はGoogle IDでユーザーを取得します
func (s *AuthService) GetUserByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	user, err := s.userRepo.GetByGoogleID(ctx, googleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
