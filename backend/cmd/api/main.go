// Command api is the Family Tree HTTP API server. It performs all dependency
// wiring (config → logger → database → repositories → services → handlers →
// router) and runs the server with graceful shutdown.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/twkakbar/family-tree/backend/internal/config"
	"github.com/twkakbar/family-tree/backend/internal/handler"
	"github.com/twkakbar/family-tree/backend/internal/repository"
	"github.com/twkakbar/family-tree/backend/internal/routes"
	"github.com/twkakbar/family-tree/backend/internal/service"
	"github.com/twkakbar/family-tree/backend/internal/validation"
	"github.com/twkakbar/family-tree/backend/pkg/password"
	"github.com/twkakbar/family-tree/backend/pkg/token"
)

// version is stamped via -ldflags "-X main.version=…" in the Docker build.
var version = "dev"

const shutdownTimeout = 15 * time.Second

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	setupLogger(cfg)
	slog.Info("starting family-tree api", "version", version, "env", cfg.AppEnv)

	if err := validation.Register(); err != nil {
		return fmt.Errorf("registering validators: %w", err)
	}

	db, err := repository.Connect(cfg.DB, cfg.LogLevel)
	if err != nil {
		return err
	}

	if cfg.RunMigrations {
		if err := repository.MigrateUp(db); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(cfg.Upload.Dir, 0o750); err != nil {
		return fmt.Errorf("creating upload dir: %w", err)
	}

	// --- Dependency wiring (constructor injection, innermost first) ---
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	jwtManager := token.NewJWTManager(cfg.Auth.JWTSecret, cfg.Auth.JWTIssuer, cfg.Auth.AccessTokenTTL)
	hasher := password.NewHasher(cfg.Auth.BcryptCost)

	authService := service.NewAuthService(userRepo, refreshTokenRepo, jwtManager, hasher, cfg.Auth.RefreshTokenTTL)
	userService := service.NewUserService(userRepo)

	if cfg.Admin.Enabled() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := authService.EnsureAdmin(ctx, cfg.Admin.Name, cfg.Admin.Email, cfg.Admin.Password); err != nil {
			cancel()
			return fmt.Errorf("bootstrapping admin: %w", err)
		}
		cancel()
	}

	router := routes.New(cfg, jwtManager, routes.Handlers{
		Auth:   handler.NewAuthHandler(authService, cfg.Auth),
		Users:  handler.NewUserHandler(userService),
		Health: handler.NewHealthHandler(db, version),
	})

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}

	// --- Serve with graceful shutdown ---
	errCh := make(chan error, 1)
	go func() {
		slog.Info("http server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("http server: %w", err)
	case sig := <-stop:
		slog.Info("shutting down", "signal", sig.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
	slog.Info("shutdown complete")
	return nil
}

// setupLogger configures the process-wide slog default: JSON in production
// (machine-ingestible), text in development (human-readable).
func setupLogger(cfg *config.Config) {
	level := slog.LevelInfo
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	var h slog.Handler
	if cfg.IsProduction() {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	} else {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}
	slog.SetDefault(slog.New(h))
}
