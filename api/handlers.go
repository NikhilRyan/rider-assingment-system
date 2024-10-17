package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"rider-assignment-system/cache"
	"rider-assignment-system/database"
	"rider-assignment-system/geohash"
	"rider-assignment-system/matching"
	"rider-assignment-system/models"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

// RequestRide handles rider's ride requests
func RequestRide(w http.ResponseWriter, r *http.Request) {
	var tripRequest struct {
		RiderID  int64   `json:"rider_id"`
		StartLat float64 `json:"start_latitude"`
		StartLon float64 `json:"start_longitude"`
		EndLat   float64 `json:"end_latitude"`
		EndLon   float64 `json:"end_longitude"`
	}

	err := json.NewDecoder(r.Body).Decode(&tripRequest)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Find the nearest available driver
	driver, err := matching.FindNearestDriver(tripRequest.StartLat, tripRequest.StartLon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Create a new trip
	var tripID int64
	err = database.DB.QueryRow(
		`INSERT INTO trips (rider_id, driver_id, start_latitude, start_longitude, end_latitude, end_longitude, status)
         VALUES ($1, $2, $3, $4, $5, $6, 'requested') RETURNING id`,
		tripRequest.RiderID, driver.ID, tripRequest.StartLat, tripRequest.StartLon, tripRequest.EndLat, tripRequest.EndLon,
	).Scan(&tripID)
	if err != nil {
		http.Error(w, "Failed to create trip", http.StatusInternalServerError)
		return
	}

	// Update driver's status to 'on_trip' in the database
	_, err = database.DB.Exec(`UPDATE drivers SET status='on_trip' WHERE id=$1`, driver.ID)
	if err != nil {
		http.Error(w, "Failed to update driver status", http.StatusInternalServerError)
		return
	}

	// Remove driver from Redis cache
	ctx := context.Background()
	driverHash := driver.Geohash
	driverJSON, _ := json.Marshal(driver)
	cache.RedisClient.SRem(ctx, fmt.Sprintf("drivers:%s", driverHash), driverJSON)

	// Respond to the rider with driver details
	response := map[string]interface{}{
		"message": "Driver assigned",
		"trip_id": tripID,
		"driver":  driver,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateDriverLocation handles updates to driver's location
func UpdateDriverLocation(w http.ResponseWriter, r *http.Request) {
	var locationUpdate struct {
		DriverID  int64   `json:"driver_id"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Status    string  `json:"status"` // Optional: "available" or "on_trip"
	}

	err := json.NewDecoder(r.Body).Decode(&locationUpdate)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get current driver data
	var currentDriver models.Driver
	err = database.DB.QueryRow(
		`SELECT id, name, latitude, longitude, geohash, status FROM drivers WHERE id=$1`,
		locationUpdate.DriverID,
	).Scan(
		&currentDriver.ID,
		&currentDriver.Name,
		&currentDriver.Latitude,
		&currentDriver.Longitude,
		&currentDriver.Geohash,
		&currentDriver.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Driver not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Calculate new geohash
	newGeohash := geohash.Encode(locationUpdate.Latitude, locationUpdate.Longitude, 5)

	// Update driver's location and status in the database
	status := locationUpdate.Status
	if status == "" {
		status = currentDriver.Status
	}
	_, err = database.DB.Exec(
		`UPDATE drivers SET latitude=$1, longitude=$2, geohash=$3, status=$4 WHERE id=$5`,
		locationUpdate.Latitude, locationUpdate.Longitude, newGeohash, status, locationUpdate.DriverID,
	)
	if err != nil {
		http.Error(w, "Failed to update driver", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	// Remove driver from old geohash set in Redis
	if currentDriver.Geohash != "" {
		currentDriverJSON, _ := json.Marshal(currentDriver)
		cache.RedisClient.SRem(ctx, fmt.Sprintf("drivers:%s", currentDriver.Geohash), currentDriverJSON)
	}

	// Add driver to new geohash set in Redis if status is 'available'
	if status == "available" {
		updatedDriver := models.Driver{
			ID:        locationUpdate.DriverID,
			Name:      currentDriver.Name,
			Latitude:  locationUpdate.Latitude,
			Longitude: locationUpdate.Longitude,
			Geohash:   newGeohash,
			Status:    status,
		}
		updatedDriverJSON, _ := json.Marshal(updatedDriver)
		cache.RedisClient.SAdd(ctx, fmt.Sprintf("drivers:%s", newGeohash), updatedDriverJSON)
	}

	// Respond with success message
	response := map[string]string{"message": "Driver location updated"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DriverStatusUpdate allows drivers to update their status (e.g., from 'on_trip' to 'available')
func DriverStatusUpdate(w http.ResponseWriter, r *http.Request) {
	var statusUpdate struct {
		DriverID int64  `json:"driver_id"`
		Status   string `json:"status"` // "available", "on_trip"
	}

	err := json.NewDecoder(r.Body).Decode(&statusUpdate)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update driver's status in the database
	_, err = database.DB.Exec(
		`UPDATE drivers SET status=$1 WHERE id=$2`,
		statusUpdate.Status, statusUpdate.DriverID,
	)
	if err != nil {
		http.Error(w, "Failed to update driver status", http.StatusInternalServerError)
		return
	}

	// Update Redis cache accordingly
	var driver models.Driver
	err = database.DB.QueryRow(
		`SELECT id, name, latitude, longitude, geohash FROM drivers WHERE id=$1`,
		statusUpdate.DriverID,
	).Scan(
		&driver.ID,
		&driver.Name,
		&driver.Latitude,
		&driver.Longitude,
		&driver.Geohash,
	)
	if err != nil {
		http.Error(w, "Failed to retrieve driver data", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	driverJSON, _ := json.Marshal(driver)
	driverKey := fmt.Sprintf("drivers:%s", driver.Geohash)

	if statusUpdate.Status == "available" {
		cache.RedisClient.SAdd(ctx, driverKey, driverJSON)
	} else {
		cache.RedisClient.SRem(ctx, driverKey, driverJSON)
	}

	// Respond with success message
	response := map[string]string{"message": "Driver status updated"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetDriver handles fetching driver details by ID
func GetDriver(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driverIDStr := vars["driver_id"]
	driverID, err := strconv.ParseInt(driverIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid driver ID", http.StatusBadRequest)
		return
	}

	var driver models.Driver
	err = database.DB.QueryRow(
		`SELECT id, name, latitude, longitude, geohash, status FROM drivers WHERE id=$1`,
		driverID,
	).Scan(
		&driver.ID,
		&driver.Name,
		&driver.Latitude,
		&driver.Longitude,
		&driver.Geohash,
		&driver.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Driver not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(driver)
}

// GetTrip handles fetching trip details by ID
func GetTrip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tripIDStr := vars["trip_id"]
	tripID, err := strconv.ParseInt(tripIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid trip ID", http.StatusBadRequest)
		return
	}

	var trip models.Trip
	err = database.DB.QueryRow(
		`SELECT id, rider_id, driver_id, start_latitude, start_longitude, end_latitude, end_longitude, status FROM trips WHERE id=$1`,
		tripID,
	).Scan(
		&trip.ID,
		&trip.RiderID,
		&trip.DriverID,
		&trip.StartLat,
		&trip.StartLon,
		&trip.EndLat,
		&trip.EndLon,
		&trip.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Trip not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trip)
}

// CreateDriver handles registering a new driver
func CreateDriver(w http.ResponseWriter, r *http.Request) {
	var driver models.Driver
	err := json.NewDecoder(r.Body).Decode(&driver)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Calculate geohash if latitude and longitude are provided
	if driver.Latitude != 0 && driver.Longitude != 0 {
		driver.Geohash = geohash.Encode(driver.Latitude, driver.Longitude, 5)
	}

	// Set default status if not provided
	if driver.Status == "" {
		driver.Status = "available"
	}

	// Insert new driver into the database
	err = database.DB.QueryRow(
		`INSERT INTO drivers (name, latitude, longitude, geohash, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		driver.Name, driver.Latitude, driver.Longitude, driver.Geohash, driver.Status,
	).Scan(&driver.ID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && strings.Contains(pgErr.Message, "duplicate key") {
			http.Error(w, "Driver already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to create driver", http.StatusInternalServerError)
		}
		return
	}

	// Add driver to Redis cache if status is 'available' and geohash is set
	if driver.Status == "available" && driver.Geohash != "" {
		ctx := context.Background()
		driverJSON, _ := json.Marshal(driver)
		cache.RedisClient.SAdd(ctx, fmt.Sprintf("drivers:%s", driver.Geohash), driverJSON)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(driver)
}

// CreateRider handles registering a new rider
func CreateRider(w http.ResponseWriter, r *http.Request) {
	var rider models.Rider
	err := json.NewDecoder(r.Body).Decode(&rider)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Insert new rider into the database
	err = database.DB.QueryRow(
		`INSERT INTO riders (name) VALUES ($1) RETURNING id`,
		rider.Name,
	).Scan(&rider.ID)
	if err != nil {
		http.Error(w, "Failed to create rider", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rider)
}

// CompleteTrip handles marking a trip as completed
func CompleteTrip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tripIDStr := vars["trip_id"]
	tripID, err := strconv.ParseInt(tripIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid trip ID", http.StatusBadRequest)
		return
	}

	// Update trip status to 'completed'
	_, err = database.DB.Exec(
		`UPDATE trips SET status='completed' WHERE id=$1`,
		tripID,
	)
	if err != nil {
		http.Error(w, "Failed to update trip", http.StatusInternalServerError)
		return
	}

	// Get the driver associated with the trip
	var driverID int64
	err = database.DB.QueryRow(
		`SELECT driver_id FROM trips WHERE id=$1`,
		tripID,
	).Scan(&driverID)
	if err != nil {
		http.Error(w, "Failed to retrieve trip details", http.StatusInternalServerError)
		return
	}

	// Update driver's status to 'available' in the database
	_, err = database.DB.Exec(
		`UPDATE drivers SET status='available' WHERE id=$1`,
		driverID,
	)
	if err != nil {
		http.Error(w, "Failed to update driver status", http.StatusInternalServerError)
		return
	}

	// Add driver back to Redis cache
	var driver models.Driver
	err = database.DB.QueryRow(
		`SELECT id, name, latitude, longitude, geohash, status FROM drivers WHERE id=$1`,
		driverID,
	).Scan(
		&driver.ID,
		&driver.Name,
		&driver.Latitude,
		&driver.Longitude,
		&driver.Geohash,
		&driver.Status,
	)
	if err != nil {
		http.Error(w, "Failed to retrieve driver data", http.StatusInternalServerError)
		return
	}

	if driver.Status == "available" && driver.Geohash != "" {
		ctx := context.Background()
		driverJSON, _ := json.Marshal(driver)
		cache.RedisClient.SAdd(ctx, fmt.Sprintf("drivers:%s", driver.Geohash), driverJSON)
	}

	response := map[string]string{"message": "Trip completed"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
