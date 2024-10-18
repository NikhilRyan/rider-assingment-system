package matching

import (
	"context"
	"encoding/json"
	"fmt"
	"rider-assignment-system/cache"
	"rider-assignment-system/geohash"
	"rider-assignment-system/models"
)

func FindNearestDriver(riderLat, riderLon float64) (*models.Driver, error) {
	riderHash := geohash.Encode(riderLat, riderLon, 5)
	neighbors := geohash.GetNeighbors(riderHash)
	neighbors = append(neighbors, riderHash)

	ctx := context.Background()

	for _, hash := range neighbors {
		drivers, err := cache.Rdb.SMembers(ctx, fmt.Sprintf("drivers:%s", hash)).Result()
		if err != nil {
			continue
		}
		for _, driverStr := range drivers {
			var driver models.Driver
			json.Unmarshal([]byte(driverStr), &driver)
			if driver.Status == "available" {
				return &driver, nil
			}
		}
	}
	return nil, fmt.Errorf("no available drivers nearby")
}
