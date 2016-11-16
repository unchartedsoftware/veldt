# Prism

>Harness the full spectrum of your data.

## Dependencies

Requires the [Go](https://golang.org/) programming language binaries with the `GOPATH` environment variable specified and `$GOPATH/bin` in your `PATH`.

## Installation

### Using `go get`:

If your project does not use the vendoring tool [Glide](https://glide.sh) to manage dependencies, you can install this package like you would any other:

```bash
go get github.com/unchartedsoftware/prism
```

While this is the simplest way to install the package, due to how `go get` resolves transitive dependencies it may result in version incompatibilities.

### Using `glide get`:

This is the recommended way to install the package and ensures all transitive dependencies are resolved to their compatible versions.

```bash
glide get github.com/unchartedsoftware/prism
```

NOTE: Requires [Glide](https://glide.sh) along with [Go](https://golang.org/) version 1.6, or version 1.5 with the `GO15VENDOREXPERIMENT` environment variable set to `1`.

## Usage

The package provides facilities to implement and connect custom tiling analytics to persistent in-memory storage services.

## Example

This minimalistic application shows how to register tile and meta data generators and connect them to a redis store.

```go
package main

import (
	"math"

	"github.com/unchartedsoftware/prism/generation/elastic"
	"github.com/unchartedsoftware/prism/generation/meta"
	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/store"
	"github.com/unchartedsoftware/prism/store/redis"
	"github.com/unchartedsoftware/prism/store/compress/gzip"
)

func GenerateMetaData(m *meta.Request) ([]byte, error) {
	// Generate meta data, this call will block until the response is ready
	// in the store.
	err := meta.GenerateMeta(m)
	if err != nil {
		return nil, err
	}
	// Retrieve the meta data form the store.
	return meta.GetMetaFromStore(m)
}

func GenerateTileData(t *tile.Request) ([]byte, error) {
	// Generate a tile, this call will block until the tile is ready in the store.
	err := prism.GenerateTile(t)
	if err != nil {
		return nil, err
	}
	// Retrieve the tile form the store.
	return tile.GetTileFromStore(t)
}

func NewTilePipeline() prism.Pipeline {
	// Create elasticsearch pipeline
	pipeline := elastic.NewPipeline("http://localhost:9200")
	// register elasticsearch tile types
	pipeline.Register("heatmap", elastic.HeatmapTile)
	pipeline.Register("top_term_count", elastic.TopTermCountTile)
	pipeline.Store()
	return pipeline
}

func NewMetaPipeline() prism.Pipeline {
	// Create elasticsearch pipeline
	pipeline := elastic.NewMetaPipeline("http://localhost:9200")
	// register elasticsearch tile types
	pipeline.Register("default", elastic.DefaultMeta)
	return pipeline
}

func main() {

	// Create pipeline

	pipeline := prism.NewPipeline()

	// Add query types to the pipeline
	pipeline.Query("exists", elastic.NewExist)
	pipeline.Query("has", elastic.NewHas)
	pipeline.Query("equals", elastic.NewEquals)
	pipeline.Query("range", elastic.NewRange)

	// Add tiles types to the pipeline
	pipeline.Tile("heatmap", elastic.NewHeatmapTile("localhost", "9200"))
	pipeline.Tile("wordcloud", elastic.NewWordcloudTile("localhost", "9200"))

	// Set the maximum concurrent tile requests
	pipeline.SetMaxConcurrent(32)
	// Set the tile requests queue length
	pipeline.SetQueueLength(1024)

	// Add meta types to the pipeline
	pipeline.Meta("default", elastic.DefaultMeta("localhost", "9200"))

	// Add a store to the pipeline
	pipeline.Store(redis.NewConnection("localhost", "6379", -1))

	// register the pipeline
	prism.Register("elastic", pipeline)

	prism.GenerateTile("elastic",
		`
		{
			"uri": "twitter_index0",
			"coord": {
				"z": 4,
				"x": 12,
				"y": 8
			},
			"tile": {
				"heatmap": {
					"xField": "pixel.x",
					"yField": "pixel.y",
					"left": 0,
					"right": 2<<32,
					"bottom": 0,
					"top": 2<<32,
					"resolution": 256
				}
			},
			"query": [
				{
					"equals": {
						"field": "name",
						"value": "john"
					}
				},
				"AND",
				{
					"range": {
						"field": "age",
						"gte": 19
					}
				}
			]
		}
		`)

	// Create a request for a `heatmap` tile.
	t := &tile.Request{
		Pipeline: "elastic",
		Tile: Heatmap{
			XField: "x",
			YField: "y",
			Left: 0,
			Right: math.Pow(2, 32),
			Bottom: 0,
			Top: math.Pow(2, 32),
			Resolution: 256,
		},
		URI: "test_index",
		Store: "redis",
		Coord: binning.Coord{
			Z: 4,
			X: 12,
			y: 12,
		},
		Query: elastic.BinaryExpression{
			Left: elastic.Equals{
				Field: "name",
				Value: "john",
			},
			Op: "AND",
			Right: elastic.Range{
				Field: "age",
				GTE: 19,
			},
		}
	}


	////


	// Create a request for `default` meta data.
	m := &meta.Request{
		Type: "elastic",
		Meta: "default",
		URI: "test_index",
	}

	// Create a request for a `heatmap` tile.
	t := &tile.Request{
		Type: "elastic",
		Tile: &elastic.Heatmap{
			XField: "x",
			YField: "y",
			Left: 0,
			Right: math.Pow(2, 32),
			Bottom: 0,
			Top: math.Pow(2, 32),
			Resolution: 256,
		},
		URI: "test_index",
		Store: "redis",
		Coord: &binning.Coord{
			Z: 4,
			X: 12,
			y: 12,
		},
		Query: &elastic.BinaryExpression{
			Left: &elastic.Equals{
				Field: "name",
				Value: "john",
			},
			Op: elastic.And,
			Right: &elastic.Range{
				Field: "age",
				GTE: 19,
			},
		}
	}

	// Generate meta data
	md, err := GenerateMetaData(m)

	// Generate tile data
	td, err := GenerateTileData(t)
}
```


## Development

Clone the repository:

```bash
mkdir -p $GOPATH/src/github.com/unchartedsoftware
cd $GOPATH/src/github.com/unchartedsoftware
git clone git@github.com:unchartedsoftware/prism.git
```

Install dependencies:

```bash
cd prism
make install
```

Debugging ES queries:

To debug queries sent to Elasticsearch by the tiles, marshall the query and aggregation sources to byte arrays, then convert to strings and print out the result.

```go
func (g *ExampleGenerator) GetTile() ([]byte, error) {
	// temp debug query
	query, err := g.getQuery().Source()
	bytes, err := json.Marshal(query)
	if err != nil {
		log.Error(err)
	}
	log.Debug(string(marshalledQuery))

	agg, err := g.getAgg().Source()
	marshalledAgg, err := json.Marshal(agg)
	if err != nil {
		log.Error(err)
	}
	log.Debug(string(marshalledAgg))
}
```
