package salt

import (
	"github.com/unchartedsoftware/veldt"
)

// Meta contains information about how to request metadata from a salt server
type Meta struct {
	rmqConfig *Config // The configuration defining how we connect to the RabbitMQ server
}

// NewMeta instantiates and returns a pointer to a new generator.
func NewMeta(rmqConfig *Config) veldt.MetaCtor {
	return func() (veldt.Meta, error) {
		return &Meta{rmqConfig}, nil
	}
}

// Create creates a metadata request
func (meta *Meta) Create(uri string) ([]byte, error) {
	connection, err := NewConnection(meta.rmqConfig)
	if err != nil {
		return nil, err
	}

	// TODO: Always transmit full dataset description and tile type with every metadata request,
	// the former in case the server has restarted, the later so that the server can return us
	// appropriate metadata
	return connection.QueryMetadata([]byte(uri))
}

// Parse gets the arguments a metadata constructor will need to create
// metadata requests.  Currently, there is no such information needed,
// so this is a no-op.
func (meta *Meta) Parse(params map[string]interface{}) error {
	return nil
}
