package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
	"log"
	"rider-assignment-system/config"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() error {
	dbHost := config.GetEnv("DB_HOST", "localhost")
	dbPort := config.GetEnv("DB_PORT", "5432")
	dbUser := config.GetEnv("DB_USER", "postgres")
	dbPassword := config.GetEnv("DB_PASSWORD", "postgres")
	dbName := config.GetEnv("DB_NAME", "matcha")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Test the connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connection established.")
	return nil
}

// GetDB returns the current database connection
func GetDB() *sql.DB {
	return DB
}

// Connect creates a new connection for readiness check
func Connect(dsn string) (*sql.DB, error) {
	return sql.Open("postgres", dsn)
}
