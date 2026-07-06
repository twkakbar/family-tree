package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

// pgUniqueViolation is the PostgreSQL error code for unique constraint violations.
const pgUniqueViolation = "23505"

// userRepository is the GORM implementation of domain.UserRepository.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository builds the PostgreSQL-backed user repository.
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if isUniqueViolation(err) {
			return domain.ErrEmailTaken
		}
		return fmt.Errorf("creating user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting user by id: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	// email is citext, so the comparison is case-insensitive in the database.
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		if isUniqueViolation(err) {
			return domain.ErrEmailTaken
		}
		return fmt.Errorf("updating user: %w", err)
	}
	return nil
}

// isUniqueViolation reports whether err is a PostgreSQL 23505 error.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation
}
