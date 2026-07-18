package config

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	AppName        string
	AppEnv         string
	Port           string
	GinMode        string
	FrontendURL    string
	UploadDir      string
	MaxUploadMB    int
	JWTSecret      string
	JWTExpireHours int
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string

	// Run embedded DB migrations on startup
	RunMigrations bool

	// CORS
	CORSAllowedOrigins []string

	// Rate limiting (token bucket, per client IP)
	RateLimitEnabled bool
	RateLimitRPS     float64
	RateLimitBurst   int
}

func LoadEnvConfig() (*EnvConfig, error) {
	_ = godotenv.Load()

	cfg := &EnvConfig{
		AppName:            getEnv("APP_NAME", "UMSRMS API"),
		AppEnv:             strings.ToLower(getEnv("APP_ENV", "development")),
		Port:               getEnv("PORT", "8080"),
		GinMode:            getEnv("GIN_MODE", "debug"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:5173"),
		UploadDir:          getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadMB:        getEnvAsInt("MAX_UPLOAD_MB", 5),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTExpireHours:     getEnvAsInt("JWT_EXPIRE_HOURS", 24),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "umsrms"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		RunMigrations:      getEnvAsBool("RUN_MIGRATIONS", true),
		CORSAllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),

		RateLimitEnabled: getEnvAsBool("RATE_LIMIT_ENABLED", true),
		RateLimitRPS:     getEnvAsFloat("RATE_LIMIT_RPS", 10),
		RateLimitBurst:   getEnvAsInt("RATE_LIMIT_BURST", 20),
	}

	if cfg.AppEnv == "production" {
		if cfg.JWTSecret == "" {
			return nil, fmt.Errorf("JWT_SECRET is required in production")
		}
	}

	return cfg, nil
}

// Address returns the HTTP bind address from configured port.
func (c *EnvConfig) Address() string {
	return ":" + c.Port
}

// PostgresDSN returns a PostgreSQL DSN string for database connection.
func (c *EnvConfig) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvAsFloat(key string, fallback float64) float64 {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvAsBool(key string, fallback bool) bool {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

// getEnvAsSlice parses a comma-separated env var into a trimmed string slice.
func getEnvAsSlice(key string, fallback []string) []string {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return fallback
	}

	return result
}

// AllowAllOrigins reports whether CORS is configured to accept any origin.
func (c *EnvConfig) AllowAllOrigins() bool {
	return slices.Contains(c.CORSAllowedOrigins, "*")
}
