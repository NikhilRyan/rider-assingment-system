package geohash

import (
	"github.com/mmcloughlin/geohash"
)

// Encode coordinates into a geohash with specified precision.
func Encode(lat, lon float64, precision uint) string {
	return geohash.EncodeWithPrecision(lat, lon, precision)
}

// GetNeighbors returns the geohashes of neighboring cells.
func GetNeighbors(hash string) []string {
	neighbors := geohash.Neighbors(hash)
	return neighbors
}
