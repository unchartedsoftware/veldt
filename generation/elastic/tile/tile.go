package tile

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"
)

// Tile represents an elasticsearch based tile generator.
type Tile struct {
	host   string
	port   string
	client *elastic.Client
}

// SetTileParams sets the params for the specific generator.
func SetTileParams(generator *Tile, host string, port string) error {
	client, err := NewClient(host, port)
	if err != nil {
		return err
	}
	generator.host = host
	generator.port = port
	generator.client = client
	return nil
}
