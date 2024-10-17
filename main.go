package main

import (
	"log"
	"net/http"
	"rider-assignment-system/api"
	"rider-assignment-system/cache"
	"rider-assignment-system/config"
	"rider-assignment-system/database"

	"github.com/gorilla/handlers"
)

func main() {
	// Initialize configuration
	config.InitConfig()

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}

	// Initialize Redis
	cache.InitRedis()

	// Register routes
	router := api.RegisterRoutes()

	// Start the server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(router)))
}
