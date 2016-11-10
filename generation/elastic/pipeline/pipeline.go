package generation

import (
	"encoding/binary"
	"math"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/param"
)

type Pipeline struct {
	Host string
	Port string
}

// NewPipeline instantiates a new elastic pipeline.
func NewPipeline(host, port, uri, typ string) tile.PipelineCtor {
	return func() (tile.Pipeline, error) {
		return &Pipeline{
			Host: host,
			Port: port,
		}, nil
	}
}

func (p *Pipeline) CreateTile(req tile.Request, query query.Query, generator tile.Generator) {
	// get client
	client, err := NewClient(p.Host, p.Port)
	if err != nil {
		return nil, err
	}

	// create search service
	search := p.client.Search().
		Index(req.URI).
		Size(0)

	// create root query
	root := elastic.NewBoolQuery()
	// apply generator to the query
	generator.ApplyQuery(root)
	// apply query
	if q != nil {
		err := q.Apply(root)
		if err != nil {
			return nil, err
		}
	}
	// add query
	search.Query(query)

	// apply generator to the query
	generator.ApplyAgg(search)

	// send query
	res, err := search.Do()
	if err != nil {
		return nil, err
	}

	return generator.ParseRes(res)
}
