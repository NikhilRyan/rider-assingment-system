package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection string
	dbURL := "postgres://postgres:postgres@db:5432/matcha?sslmode=disable"

	// Wait for the database to be ready
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dbURL)
		if err == nil && db.Ping() == nil {
			break
		}
		log.Println("Waiting for the database to be ready...")
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	db.Close()

	// Run migrations
	migrationsPath := "file://database/migrations"
	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Could not start migrations: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Migrations applied successfully!")
}
