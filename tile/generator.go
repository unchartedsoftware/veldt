package tile

// Generator represents an interface for generating tile data.
type Generator interface {
	GetTile() ([]byte, error)
	GetParams() []Param
	GetHash() string
}

// GeneratorConstructor represents a function to instantiate a new generator
// from a tile request.
type GeneratorConstructor func(*Request) (Generator, error)
