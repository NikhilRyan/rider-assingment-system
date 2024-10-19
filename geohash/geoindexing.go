package geohash

import (
	"errors"
)

type GeoIndexingTechnique string

const (
	GeohashingTechnique GeoIndexingTechnique = "geohashing"
	RTreeTechnique      GeoIndexingTechnique = "rtree"
	QuadtreeTechnique   GeoIndexingTechnique = "quadtree"
)

var defaultTechnique = GeohashingTechnique
var quadtreeInstance *Quadtree

// SetDefaultTechnique sets the default geo-indexing technique
func SetDefaultTechnique(technique GeoIndexingTechnique) {
	defaultTechnique = technique
}

// InitializeGlobalQuadtree initializes a global Quadtree instance
func InitializeGlobalQuadtree(bounds Bounds) {
	quadtreeInstance = InitializeQuadtree(bounds)
}

// SearchNearbyWithRetries tries to find nearby points with a retry mechanism
func SearchNearbyWithRetries(lat, lon float64, technique GeoIndexingTechnique, maxRetries int) ([]interface{}, error) {
	if technique == "" {
		technique = defaultTechnique
	}

	radius := 1.0 // Initial search radius
	var results []interface{}

	for i := 0; i < maxRetries; i++ {
		switch technique {
		case GeohashingTechnique:
			geohash := Encode(lat, lon, 12)
			neighbors := GetNeighbors(geohash)
			for _, neighbor := range neighbors {
				results = append(results, neighbor)
			}
		case RTreeTechnique:
			rtreeResults := SearchNearbyInRTree(lat, lon, radius)
			for _, item := range rtreeResults {
				results = append(results, item)
			}
		case QuadtreeTechnique:
			if quadtreeInstance != nil {
				qtResults := quadtreeInstance.SearchNearbyInQuadtree(Point{lat, lon}, radius)
				for _, point := range qtResults {
					results = append(results, point)
				}
			}
		default:
			return nil, errors.New("unsupported geo-indexing technique")
		}

		if len(results) > 0 {
			break // If results are found, exit the loop
		}

		radius *= 2 // Increase the search radius for the next retry
	}

	if len(results) == 0 {
		return nil, errors.New("no nearby points found after maximum retries")
	}

	return results, nil
}
