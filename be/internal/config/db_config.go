package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewPostgresDB opens and validates a PostgreSQL database connection.
func NewPostgresDB(cfg *EnvConfig) (*sql.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	db, err := sql.Open("pgx", cfg.PostgresDSN())
	if err != nil {
		return nil, err
	}

	// Conservative defaults that work well for local development.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
