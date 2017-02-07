package geometry

// Coord represents a point.
type Coord struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// NewCoord instantiates and returns a pointer to a Coord.
func NewCoord(x, y float64) *Coord {
	return &Coord{
		X: x,
		Y: y,
	}
}
