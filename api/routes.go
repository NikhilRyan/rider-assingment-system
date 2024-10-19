package api

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func RegisterRoutes() http.Handler {
	router := mux.NewRouter()

	// Rider endpoints
	router.HandleFunc("/riders", CreateRider).Methods("POST")

	// Driver endpoints
	router.HandleFunc("/drivers", CreateDriver).Methods("POST")
	router.HandleFunc("/drivers/{driver_id}", GetDriver).Methods("GET")
	router.HandleFunc("/drivers/{driver_id}/status", DriverStatusUpdate).Methods("PUT")
	router.HandleFunc("/drivers/{driver_id}/location", UpdateDriverLocation).Methods("PUT")

	// Trip endpoints
	router.HandleFunc("/trips", RequestRide).Methods("POST")
	router.HandleFunc("/trips/{trip_id}", GetTrip).Methods("GET")
	router.HandleFunc("/trips/{trip_id}/complete", CompleteTrip).Methods("PUT")

	// Distance endpoint
	router.HandleFunc("/distance", DistanceHandler).Methods("POST")

	// Add the GeoIndexingHandler route
	router.HandleFunc("/geoindex", GeoIndexingHandler).Methods("GET")

	// Add CORS support
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	return cors(router)
}
