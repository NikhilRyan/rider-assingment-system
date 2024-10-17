package models

type Driver struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Geohash   string  `json:"geohash"`
	Status    string  `json:"status"` // "available", "on_trip"
}
