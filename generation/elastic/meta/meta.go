package meta

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	es "github.com/unchartedsoftware/prism/generation/elastic"
)

// Meta represents an elasticsearch based meta generator.
type Meta struct {
	host   string
	port   string
	client *elastic.Client
}

// SetMetaParams sets the params for the specific generator.
func SetMetaParams(arg interface{}, host string, port string) error {
	generator, ok := arg.(*Meta)
	if !ok {
		return fmt.Errorf("`%v` is not of type `*Meta`", arg)
	}
	client, err := es.NewClient(host, port)
	if err != nil {
		return err
	}
	generator.host = host
	generator.port = port
	generator.client = client
	return nil
}
