package migration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// RunMigrations runs the database migrations
func RunMigrations() error {
	// Database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Retry connecting to the database to ensure it's ready
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil && db.Ping() == nil {
			log.Println("Connected to the database successfully.")
			break
		}
		log.Printf("Waiting for the database to be ready... (attempt %d)", i+1)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("could not connect to the database: %v", err)
	}
	db.Close()

	// Run the migrations
	migrationsPath := "file://database/migrations"
	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		return fmt.Errorf("could not start migrations: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %v", err)
	}

	log.Println("Migrations applied successfully!")
	return nil
}
