package main

import (
	"log"
	"net/http"
	"os"
	"rider-assignment-system/geohash"
	"time"

	"rider-assignment-system/api"
	"rider-assignment-system/cache"
	"rider-assignment-system/config"
	"rider-assignment-system/database"

	"github.com/gorilla/handlers"
)

func main() {
	// Initialize configuration
	config.InitConfig()

	// Wait for the database to be ready with a retry mechanism
	if err := waitForDatabase(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Initialize the database connection
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}

	// Initialize Redis
	if err := cache.InitializeRedis(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// Set the default geo-indexing technique
	geohash.SetDefaultTechnique(geohash.GeohashingTechnique)

	// Initialize the R-tree
	geohash.InitializeRTree()

	// Initialize the Quadtree with specified bounds
	quadtreeBounds := geohash.Bounds{
		MinX: -180, MinY: -90, // Minimum bounds for latitude and longitude
		MaxX: 180, MaxY: 90, // Maximum bounds for latitude and longitude
	}
	geohash.InitializeGlobalQuadtree(quadtreeBounds)

	// Register routes for the API
	router := api.RegisterRoutes()

	// Start the HTTP server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(router)))
}

// waitForDatabase attempts to connect to the database with a retry mechanism.
func waitForDatabase() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"

	// Retry mechanism to wait for the database to be ready
	for i := 0; i < 10; i++ {
		db, err := database.Connect(dsn)
		if err == nil && db.Ping() == nil {
			log.Println("Connected to the database successfully.")
			db.Close()
			return nil
		}
		log.Printf("Waiting for the database to be ready... (attempt %d)", i+1)
		time.Sleep(3 * time.Second)
	}
	return log.Output(0, "Failed to connect to the database after multiple attempts")
}
