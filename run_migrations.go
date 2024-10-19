package main

import (
	"log"
	"rider-assignment-system/migration"
)

func main() {
	// Run the migrations
	if err := migration.RunMigrations(); err != nil {
		log.Fatalf("Migration error: %v", err)
	}
}
