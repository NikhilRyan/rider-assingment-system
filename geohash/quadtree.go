package geohash

import (
	"math"
	"sync"
)

// Point represents a point in 2D space
type Point struct {
	X, Y float64
}

// Bounds represents the boundaries of a region
type Bounds struct {
	MinX, MinY, MaxX, MaxY float64
}

// QuadtreeNode represents a node in the quadtree
type QuadtreeNode struct {
	Bounds   Bounds
	Points   []Point
	Children [4]*QuadtreeNode
}

// Quadtree represents the quadtree structure
type Quadtree struct {
	Root *QuadtreeNode
	Lock sync.Mutex
}

// InitializeQuadtree initializes a new Quadtree with given bounds
func InitializeQuadtree(bounds Bounds) *Quadtree {
	return &Quadtree{
		Root: &QuadtreeNode{Bounds: bounds},
	}
}

// Insert adds a point to the Quadtree
func (qt *Quadtree) Insert(point Point) {
	qt.Lock.Lock()
	defer qt.Lock.Unlock()
	qt.Root.insert(point)
}

// insert adds a point to a QuadtreeNode, creating children nodes if necessary
func (node *QuadtreeNode) insert(point Point) {
	if !node.contains(point) {
		return
	}
	if len(node.Points) < 4 && node.Children[0] == nil {
		node.Points = append(node.Points, point)
		return
	}
	if node.Children[0] == nil {
		node.subdivide()
	}
	for i := 0; i < 4; i++ {
		node.Children[i].insert(point)
	}
}

// contains checks if the point is within the node's bounds
func (node *QuadtreeNode) contains(point Point) bool {
	return point.X >= node.Bounds.MinX && point.X <= node.Bounds.MaxX &&
		point.Y >= node.Bounds.MinY && point.Y <= node.Bounds.MaxY
}

// subdivide splits the node into four child nodes
func (node *QuadtreeNode) subdivide() {
	midX := (node.Bounds.MinX + node.Bounds.MaxX) / 2
	midY := (node.Bounds.MinY + node.Bounds.MaxY) / 2
	node.Children[0] = &QuadtreeNode{Bounds: Bounds{node.Bounds.MinX, node.Bounds.MinY, midX, midY}}
	node.Children[1] = &QuadtreeNode{Bounds: Bounds{midX, node.Bounds.MinY, node.Bounds.MaxX, midY}}
	node.Children[2] = &QuadtreeNode{Bounds: Bounds{node.Bounds.MinX, midY, midX, node.Bounds.MaxY}}
	node.Children[3] = &QuadtreeNode{Bounds: Bounds{midX, midY, node.Bounds.MaxX, node.Bounds.MaxY}}
}

// SearchNearbyInQuadtree searches for nearby points within a given radius
func (qt *Quadtree) SearchNearbyInQuadtree(center Point, radius float64) []Point {
	qt.Lock.Lock()
	defer qt.Lock.Unlock()
	return qt.Root.searchNearby(center, radius)
}

// searchNearby finds points within a radius in a QuadtreeNode
func (node *QuadtreeNode) searchNearby(center Point, radius float64) []Point {
	if !node.intersectsCircle(center, radius) {
		return nil
	}
	var result []Point
	for _, point := range node.Points {
		if distance(point, center) <= radius {
			result = append(result, point)
		}
	}
	if node.Children[0] != nil {
		for i := 0; i < 4; i++ {
			result = append(result, node.Children[i].searchNearby(center, radius)...)
		}
	}
	return result
}

// intersectsCircle checks if a circle intersects with the node's bounds
func (node *QuadtreeNode) intersectsCircle(center Point, radius float64) bool {
	closestX := math.Max(node.Bounds.MinX, math.Min(center.X, node.Bounds.MaxX))
	closestY := math.Max(node.Bounds.MinY, math.Min(center.Y, node.Bounds.MaxY))
	dx := closestX - center.X
	dy := closestY - center.Y
	return (dx*dx + dy*dy) <= (radius * radius)
}

// distance calculates the Euclidean distance between two points
func distance(a, b Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}
