# Prism

>Harness the full spectrum of your data.

## Dependencies

Requires the [Go](https://golang.org/) programming language binaries with the `GOPATH` environment variable specified.

## Installation

```bash
go get github.com/unchartedsoftware/prism
```

## Usage

The package provides facilities to implement and connect custom tiling and analytics to persistent in-memory storage services.

## Example

This minimalistic application shows how to register tile and meta data generators and connect them to a redis store.

```go
package main

import (
    "github.com/unchartedsoftware/prism/generation/elastic"
	"github.com/unchartedsoftware/prism/generation/meta"
    "github.com/unchartedsoftware/prism/generation/tile"
    "github.com/unchartedsoftware/prism/store"
    "github.com/unchartedsoftware/prism/store/redis"
)

func GenerateMetaData(m *meta.Request, s *store.Request) ([]byte, error) {
    // Generate meta data, this call will block until the response is ready
    // in the store.
    err := meta.GenerateMeta(m, s)
    if err != nil {
    	return nil, err
    }
    // Retrieve the meta data form the store.
    return meta.GetMetaFromStore(m, s)
}

func GenerateTileData(t *tile.Request, s *store.Request) ([]byte, error) {
    // Generate a tile, this call will block until the tile is ready in the store.
    err := tile.GenerateTile(t, s)
    if err != nil {
    	return nil, err
    }
    // Retrieve the tile form the store.
    return tile.GetTileFromStore(t, s)
}

func main() {    
    // Register the in-memory store to use the redis implementation.
    store.Register("redis", redis.NewConnection)

    // Register tile and meta data generators
    tile.Register("heatmap", elastic.NewHeatmapTile)
    meta.Register("default", elastic.NewDefaultMeta)

    // Create a request for a `heatmap` tile.
    m := &meta.Request{
        Type: "default",
        Endpoint: "http://localhost:9200",
        Index: "test_index",
    }

    // Create a request for a `heatmap` tile.
    t := &tile.Request{
        Type: "heatmap",
        Endpoint: "http://localhost:9200",
        Index: "test_index",
        Coord: &binning.TileCoord{
            Z: 4,
            X: 12,
            y: 12,
        }
    }

    // Create the request to store the generated tile in the `redis` store.
    s := &store.Request{
        Type: "redis",
        Endpoint: "localhost:6379",
    }

    // Generate meta data
    md, err := GenerateMetaData(m, s)

    // Generate tile data
    td, err := GenerateTileData(t, s)
}
```
