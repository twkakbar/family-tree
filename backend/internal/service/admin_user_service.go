package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
	"github.com/twkakbar/family-tree/backend/pkg/token"
)

type adminUserService struct {
	users          domain.UserRepository
	refreshTokens  domain.RefreshTokenRepository
	passwordResets domain.PasswordResetTokenRepository
	resetTokenTTL  time.Duration
}

// NewAdminUserService wires super-admin account management.
func NewAdminUserService(
	users domain.UserRepository,
	refreshTokens domain.RefreshTokenRepository,
	passwordResets domain.PasswordResetTokenRepository,
) domain.AdminUserService {
	return &adminUserService{
		users:          users,
		refreshTokens:  refreshTokens,
		passwordResets: passwordResets,
		resetTokenTTL:  time.Hour,
	}
}

func (s *adminUserService) List(ctx context.Context, adminID uuid.UUID, search string, page, perPage int) ([]domain.User, int64, error) {
	if err := s.requireAdmin(ctx, adminID); err != nil {
		return nil, 0, err
	}
	params := normalizePage(page, perPage)
	return s.users.List(ctx, strings.TrimSpace(search), params.offset(), params.limit())
}

func (s *adminUserService) Update(ctx context.Context, adminID, userID uuid.UUID, in domain.AdminUpdateUserInput) (*domain.User, error) {
	if err := s.requireAdmin(ctx, adminID); err != nil {
		return nil, err
	}
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		name := strings.TrimSpace(*in.Name)
		if len(name) < 2 || len(name) > 100 {
			return nil, domain.ErrUnsupportedFormat
		}
		user.Name = name
	}
	if in.Role != nil {
		role := *in.Role
		if role != domain.RoleAdmin && role != domain.RoleUser {
			return nil, domain.ErrUnsupportedFormat
		}
		if user.ID == adminID && role != domain.RoleAdmin {
			return nil, domain.ErrForbidden
		}
		user.Role = role
	}
	if in.EmailVerified != nil {
		if *in.EmailVerified {
			now := time.Now().UTC()
			user.EmailVerifiedAt = &now
		} else {
			user.EmailVerifiedAt = nil
		}
	}
	user.UpdatedAt = time.Now().UTC()
	if err := s.users.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *adminUserService) Delete(ctx context.Context, adminID, userID uuid.UUID) error {
	if err := s.requireAdmin(ctx, adminID); err != nil {
		return err
	}
	if adminID == userID {
		return domain.ErrForbidden
	}
	if err := s.users.SoftDelete(ctx, userID); err != nil {
		return err
	}
	_ = s.refreshTokens.RevokeAllForUser(ctx, userID)
	_ = s.passwordResets.RevokeAllForUser(ctx, userID)
	return nil
}

func (s *adminUserService) IssuePasswordReset(ctx context.Context, adminID, userID uuid.UUID) (string, error) {
	if err := s.requireAdmin(ctx, adminID); err != nil {
		return "", err
	}
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if err := s.passwordResets.RevokeAllForUser(ctx, user.ID); err != nil {
		return "", err
	}
	raw, err := token.NewOpaque()
	if err != nil {
		return "", err
	}
	now := time.Now().UTC()
	record := &domain.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: token.HashOpaque(raw),
		ExpiresAt: now.Add(s.resetTokenTTL),
		CreatedAt: now,
	}
	if err := s.passwordResets.Create(ctx, record); err != nil {
		return "", err
	}
	return raw, nil
}

func (s *adminUserService) requireAdmin(ctx context.Context, adminID uuid.UUID) error {
	admin, err := s.users.GetByID(ctx, adminID)
	if err != nil {
		return err
	}
	if !admin.IsAdmin() {
		return domain.ErrForbidden
	}
	return nil
}
