package veldt

// Store represents an interface for connecting to, setting, and retrieving
// values from a key-value database or in-memory storage server.
type Store interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Exists(string) (bool, error)
	Close()
}

// StoreCtor represents a function that instantiates and returns a new storage
// type.
type StoreCtor func() (Store, error)
