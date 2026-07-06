package repository

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/gorm"

	"github.com/twkakbar/family-tree/backend/migrations"
)

// NewMigrator builds a golang-migrate instance backed by the SQL files
// embedded in the migrations package, reusing the GORM connection.
func NewMigrator(db *gorm.DB) (*migrate.Migrate, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("acquiring sql.DB for migrations: %w", err)
	}

	source, err := iofs.New(migrations.FS(), ".")
	if err != nil {
		return nil, fmt.Errorf("loading embedded migrations: %w", err)
	}

	driver, err := migratepg.WithInstance(sqlDB, &migratepg.Config{})
	if err != nil {
		return nil, fmt.Errorf("initializing migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("initializing migrator: %w", err)
	}
	return m, nil
}

// MigrateUp applies all pending migrations. Safe to run on every startup:
// it is a no-op when the schema is current. Intended for single-instance
// deployments; multi-instance rollouts should run cmd/migrate separately.
func MigrateUp(db *gorm.DB) error {
	m, err := NewMigrator(db)
	if err != nil {
		return err
	}

	version, dirty, _ := m.Version()
	if dirty {
		return fmt.Errorf("database is dirty at version %d; resolve manually (see docs/DATABASE.md)", version)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("migrations up to date", "version", version)
			return nil
		}
		return fmt.Errorf("applying migrations: %w", err)
	}

	newVersion, _, _ := m.Version()
	slog.Info("migrations applied", "from", version, "to", newVersion)
	return nil
}
