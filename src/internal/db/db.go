package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pkgzx/liliApi/src/pkg/config"
)

type DB struct {
	*sql.DB
}

func NewConnection(cfg *config.DatabaseConfig) (*DB, error) {
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Database connection established")
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
