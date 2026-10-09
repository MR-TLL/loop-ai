package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"backend/config"

	_ "github.com/lib/pq"
	_ "embed"
)

//go:embed init.sql
var initSQL string

func Init(cfg config.DatabaseConfig) *sql.DB {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeMins) * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	if _, err := db.Exec(initSQL); err != nil {
		log.Fatalf("Failed to run init SQL: %v", err)
	}

	return db
}

func BeginTx(ctx context.Context, db *sql.DB) (*sql.Tx, error) {
	return db.BeginTx(ctx, nil)
}
