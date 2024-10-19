package geohash

import (
	"github.com/dhconnelly/rtreego"
	"sync"
)

// SpatialPoint wraps a point to satisfy the rtreego.Spatial interface
type SpatialPoint struct {
	rtreego.Point
}

// BoundingBox returns a rectangle representing the spatial bounds of the point
func (p SpatialPoint) BoundingBox() *rtreego.Rect {
	// Create a small bounding box around the point
	zeroDistance := 0.0001 // A very small distance to represent the bounding box
	rect := p.Point.ToRect(zeroDistance)
	return &rect
}

var rtree *rtreego.Rtree
var rtreeLock sync.Mutex

// InitializeRTree initializes the R-tree for spatial indexing
func InitializeRTree() {
	rtree = rtreego.NewTree(2, 25, 50)
}

//// AddPointToRTree adds a point to the R-tree
//func AddPointToRTree(lat, lon float64) {
//	rtreeLock.Lock()
//	defer rtreeLock.Unlock()
//	point := SpatialPoint{rtreego.Point{lat, lon}}
//	rtree.Insert(point)
//}

// SearchNearbyInRTree searches for nearby points within a given radius
func SearchNearbyInRTree(lat, lon, radius float64) []rtreego.Spatial {
	rtreeLock.Lock()
	defer rtreeLock.Unlock()
	point := rtreego.Point{lat, lon}
	return rtree.SearchIntersect(point.ToRect(radius))
}
