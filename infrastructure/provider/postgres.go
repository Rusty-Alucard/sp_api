// Package provider provides database connection and configuration utilities.
package provider

import (
	"context"
	"database/sql"
	"os"
	"time"

	// PostgreSQL driver
	_ "github.com/lib/pq"
)

const (
	defaultTimeout = 5 * time.Second
)

// Connect establishes a connection to the PostgreSQL database.
func Connect() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "sp_api"
	}

	dsn := "host=" + host + " port=" + port + " user=" + user + 
		" password=" + password + " dbname=" + dbname + " sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
