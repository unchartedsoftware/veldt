# Prism

>Harness the full spectrum of your data.

## Dependencies

Requires the [Go](https://golang.org/) programming language binaries with the `GOPATH` environment variable specified.

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

## Development

Clone the repository:

```bash
mkdir $GOPATH/src/github.com/unchartedsoftware
cd $GOPATH/src/github.com/unchartedsoftware
git clone git@github.com:unchartedsoftware/prism.git
```

Install dependencies

```bash
cd prism
make deps
```

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
    err := tile.GenerateTile(t)
    if err != nil {
        return nil, err
    }
    // Retrieve the tile form the store.
    return tile.GetTileFromStore(t)
}

func main() {
    // Register the in-memory store to use the redis implementation.
    store.Register("redis", redis.NewConnection("localhost", "6379", 3600))    
    // Use gzip compression when setting / getting from the store
    store.Use(gzip.NewCompressor())

    // Register meta data generator
    meta.Register("default", elastic.NewDefaultMeta("http://localhost", "9200"))

    // Register tile data generator
    tile.Register("heatmap", elastic.NewHeatmapTile("http://localhost", "9200"))
    // Set the maximum concurrent tile requests
    tile.SetMaxConcurrent(32)
    // Set the tile requests queue length
    tile.SetQueueLength(1024)

    // Create a request for `default` meta data.
    m := &meta.Request{
        Type: "default",
        URI: "test_index",
    }

    // Create a request for a `heatmap` tile.
    t := &tile.Request{
        Type: "heatmap",
        URI: "test_index",
        Store: "redis"
        Coord: &binning.TileCoord{
            Z: 4,
            X: 12,
            y: 12,
        },
        Params: map[string]interface{}{
            "binning": map[string]interface{}{
                "x": "xField",
                "y": "yField",
                "left": 0,
                "right": math.Pow(2, 32),
                "bottom": 0,
                "top": math.Pow(2, 32),
                "resolution": 256,
            }
        }
    }

    // Generate meta data
    md, err := GenerateMetaData(m)

    // Generate tile data
    td, err := GenerateTileData(t)
}
```
