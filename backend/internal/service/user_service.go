package service

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

// userService implements domain.UserService.
type userService struct {
	users domain.UserRepository
}

// NewUserService wires the user profile business logic.
func NewUserService(users domain.UserRepository) domain.UserService {
	return &userService{users: users}
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *userService) UpdateProfile(ctx context.Context, id uuid.UUID, in domain.UpdateProfileInput) (*domain.User, error) {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Name = strings.TrimSpace(in.Name)
	if err := s.users.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
