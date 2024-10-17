package models

type Trip struct {
	ID       int64   `json:"id"`
	RiderID  int64   `json:"rider_id"`
	DriverID int64   `json:"driver_id"`
	StartLat float64 `json:"start_latitude"`
	StartLon float64 `json:"start_longitude"`
	EndLat   float64 `json:"end_latitude"`
	EndLon   float64 `json:"end_longitude"`
	Status   string  `json:"status"` // "requested", "accepted", "completed"
}
