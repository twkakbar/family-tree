// Package repository contains the persistence layer: the PostgreSQL
// connection, the migration runner and GORM implementations of the
// repository interfaces declared in the domain package. It translates all
// driver-level errors into domain errors; business logic is forbidden here.
package repository

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/twkakbar/family-tree/backend/internal/config"
)

// Connect opens the PostgreSQL connection pool with the configured limits.
func Connect(cfg config.DBConfig, logLevel string) (*gorm.DB, error) {
	gormLevel := gormlogger.Warn
	if logLevel == "debug" {
		gormLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger:  gormlogger.Default.LogMode(gormLevel),
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("acquiring sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	slog.Info("database connected",
		"host", cfg.Host, "port", cfg.Port, "database", cfg.Name)
	return db, nil
}
