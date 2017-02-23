package salt

import (
	"github.com/unchartedsoftware/veldt"
)

// DefaultMeta contains information about how to request metadata from a salt server
type DefaultMeta struct {
	rmqConfig *Configuration // The configuration defining how we connect to the RabbitMQ server
}

// NewDefaultMeta instantiates and returns a pointer to a new generator.
func NewDefaultMeta(rmqConfig *Configuration) veldt.MetaCtor {
	return func() (veldt.Meta, error) {
		return &DefaultMeta{rmqConfig}, nil
	}
}

// Create creates a metadata request
func (meta *DefaultMeta) Create (uri string) ([]byte, error) {
	connection, err := NewConnection(meta.rmqConfig)
	if err != nil {
		return nil, err
	}

	qName := meta.rmqConfig.serverQueue
	result, err := connection.Query(qName, []byte("test metadata query"))
	return result, err
}

// Parse does something with a metadata request, but I've no idea what
func (meta *DefaultMeta) Parse (params map[string]interface{}) error {
	return nil
}

