package meta

// Generator represents an interface for generating meta data.
type Generator interface {
	GetMeta(string) ([]byte, error)
}

// GeneratorConstructor represents a function to instantiate a new generator
// from a meta data request.
type GeneratorConstructor func(*Request) (Generator, error)
