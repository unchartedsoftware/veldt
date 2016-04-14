package store

var (
	// the bound compressor
	compressor Compressor
	// default to no compression
	disabled = true
)

// Compressor represents an interface for compressing and decompressing the
// generated tile data before adding and after retrieving it from the store.
type Compressor interface {
	Compress([]byte) ([]byte, error)
	Decompress([]byte) ([]byte, error)
}

// CompressorConstructor represents a function to instantiate a new Compressor.
type CompressorConstructor func() Compressor

// Use registers the provided compressor to be used.
func Use(comp CompressorConstructor) {
	compressor = comp()
	disabled = false
}

func compress(data []byte) ([]byte, error) {
	if disabled {
		return data, nil
	}
	return compressor.Compress(data)
}

func decompress(data []byte) ([]byte, error) {
	if disabled {
		return data, nil
	}
	return compressor.Decompress(data)
}
