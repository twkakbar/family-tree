// Package config loads and validates all application configuration from
// environment variables. It is the ONLY place that reads the environment;
// every other package receives typed values via dependency injection.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Environment names accepted in APP_ENV.
const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
	EnvTest        = "test"
)

// Config is the complete, validated application configuration.
type Config struct {
	AppEnv   string
	HTTPPort int
	LogLevel string

	DB             DBConfig
	Auth           AuthConfig
	CORS           CORSConfig
	Upload         UploadConfig
	Admin          AdminBootstrap
	RoyalFamily    RoyalFamilyConfig
	SMTP           SMTPConfig
	PublicURL      string
	RunMigrations  bool
	TrustedProxies []string
}

// DBConfig holds PostgreSQL connection settings.
type DBConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DSN renders the connection string used by GORM and golang-migrate.
func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

// URL renders the postgres:// URL form of the DSN.
func (d DBConfig) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

// AuthConfig holds token and password-hashing settings.
type AuthConfig struct {
	JWTSecret       string
	JWTIssuer       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	BcryptCost      int
	CookieSecure    bool
	CookieDomain    string
}

// CORSConfig holds cross-origin settings (relevant for hybrid development;
// production is same-origin behind nginx).
type CORSConfig struct {
	AllowedOrigins []string
}

// UploadConfig holds file upload settings.
type UploadConfig struct {
	Dir         string
	MaxUploadMB int64
}

// AdminBootstrap optionally seeds an admin account on startup.
type AdminBootstrap struct {
	Email    string
	Password string
	Name     string
}

// Enabled reports whether bootstrap credentials are configured.
func (a AdminBootstrap) Enabled() bool { return a.Email != "" && a.Password != "" }

// RoyalFamilyConfig names the singleton public family seeded at startup.
type RoyalFamilyConfig struct {
	Name        string
	Description string
}

// SMTPConfig is outbound mail. Empty host disables verification email.
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// Load reads configuration from the environment. In non-production it first
// loads a local .env file when present (developer convenience; production
// containers receive real environment variables from compose).
func Load() (*Config, error) {
	if os.Getenv("APP_ENV") != EnvProduction {
		// Ignore absence: .env is optional.
		_ = godotenv.Load()
	}

	cfg := &Config{
		AppEnv:   getEnv("APP_ENV", EnvDevelopment),
		HTTPPort: getEnvInt("HTTP_PORT", 8080),
		LogLevel: strings.ToLower(getEnv("LOG_LEVEL", "info")),
		DB: DBConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "familytree"),
			Password:        getEnv("DB_PASSWORD", ""),
			Name:            getEnv("DB_NAME", "familytree"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
		},
		Auth: AuthConfig{
			JWTSecret:       getEnv("JWT_SECRET", ""),
			JWTIssuer:       getEnv("JWT_ISSUER", "family-tree-api"),
			AccessTokenTTL:  getEnvDuration("ACCESS_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL: getEnvDuration("REFRESH_TOKEN_TTL", 720*time.Hour),
			BcryptCost:      getEnvInt("BCRYPT_COST", 12),
			CookieSecure:    getEnvBool("COOKIE_SECURE", false),
			CookieDomain:    getEnv("COOKIE_DOMAIN", ""),
		},
		CORS: CORSConfig{
			AllowedOrigins: splitNonEmpty(getEnv("CORS_ALLOWED_ORIGINS", "")),
		},
		Upload: UploadConfig{
			Dir:         getEnv("UPLOAD_DIR", "./uploads"),
			MaxUploadMB: int64(getEnvInt("MAX_UPLOAD_SIZE_MB", 10)),
		},
		Admin: AdminBootstrap{
			Email:    getEnv("ADMIN_EMAIL", ""),
			Password: getEnv("ADMIN_PASSWORD", ""),
			Name:     getEnv("ADMIN_NAME", "Administrator"),
		},
		RoyalFamily: RoyalFamilyConfig{
			Name:        getEnv("ROYAL_FAMILY_NAME", "Kesultanan Aceh"),
			Description: getEnv("ROYAL_FAMILY_DESCRIPTION", "Silsilah keluarga kerajaan Aceh untuk dijelajahi publik."),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", ""),
			Port:     getEnvInt("SMTP_PORT", 465),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", ""),
		},
		PublicURL:      strings.TrimRight(getEnv("APP_PUBLIC_URL", ""), "/"),
		RunMigrations:  getEnvBool("RUN_MIGRATIONS", true),
		TrustedProxies: splitNonEmpty(getEnv("TRUSTED_PROXIES", "")),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// IsProduction reports whether the app runs with production settings.
func (c *Config) IsProduction() bool { return c.AppEnv == EnvProduction }

func (c *Config) validate() error {
	var errs []error

	switch c.AppEnv {
	case EnvDevelopment, EnvProduction, EnvTest:
	default:
		errs = append(errs, fmt.Errorf("APP_ENV must be one of development|production|test, got %q", c.AppEnv))
	}
	if len(c.Auth.JWTSecret) < 32 {
		errs = append(errs, errors.New("JWT_SECRET is required and must be at least 32 characters"))
	}
	if c.DB.Password == "" && c.AppEnv == EnvProduction {
		errs = append(errs, errors.New("DB_PASSWORD is required in production"))
	}
	if c.Auth.AccessTokenTTL <= 0 || c.Auth.RefreshTokenTTL <= 0 {
		errs = append(errs, errors.New("token TTLs must be positive durations"))
	}
	smtpOn := c.SMTP.Host != "" || c.SMTP.Username != "" || c.SMTP.Password != "" || c.SMTP.From != ""
	if smtpOn && (c.SMTP.Host == "" || c.SMTP.Username == "" || c.SMTP.Password == "" || c.SMTP.From == "" || c.PublicURL == "") {
		errs = append(errs, errors.New("SMTP_HOST, SMTP_USERNAME, SMTP_PASSWORD, SMTP_FROM, and APP_PUBLIC_URL are all required when SMTP is configured"))
	}
	if c.IsProduction() && !c.Auth.CookieSecure {
		// Deliberately a warning-level decision: HTTP-only deployments exist
		// during initial VM setup. DEPLOYMENT.md instructs enabling it with TLS.
		fmt.Fprintln(os.Stderr, "WARNING: COOKIE_SECURE=false in production - enable it once HTTPS is configured")
	}

	return errors.Join(errs...)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v)
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(strings.TrimSpace(v)); err == nil {
			return b
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok {
		if d, err := time.ParseDuration(strings.TrimSpace(v)); err == nil {
			return d
		}
	}
	return fallback
}

func splitNonEmpty(csv string) []string {
	var out []string
	for _, part := range strings.Split(csv, ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}
