package database

import (
	"database/sql"
	"fmt"
	"log"
	"rider-assignment-system/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	cfg := config.Cfg.DB
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	DB = db
	log.Println("Database connected.")
	return nil
}
