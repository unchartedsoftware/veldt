<img width="600" src="https://rawgit.com/unchartedsoftware/veldt/master/logo.svg" alt="veldt" />

> Scalable on-demand tile-based analytics

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](http://godoc.org/github.com/unchartedsoftware/veldt)
[![Build Status](https://travis-ci.org/unchartedsoftware/veldt.svg?branch=master)](https://travis-ci.org/unchartedsoftware/veldt)
[![Go Report Card](https://goreportcard.com/badge/github.com/unchartedsoftware/veldt)](https://goreportcard.com/report/github.com/unchartedsoftware/veldt)

## Dependencies

Requires the [Go](https://golang.org/) programming language binaries with the `GOPATH` environment variable specified and `$GOPATH/bin` in your `PATH`.

## Installation

### Using `go get`:

If your project does not use the vendoring tool [Glide](https://glide.sh) to manage dependencies, you can install this package like you would any other:

```bash
go get github.com/unchartedsoftware/veldt
```

While this is the simplest way to install the package, due to how `go get` resolves transitive dependencies it may result in version incompatibilities.

### Using `glide get`:

This is the recommended way to install the package and ensures all transitive dependencies are resolved to their compatible versions.

```bash
glide get github.com/unchartedsoftware/veldt
```

NOTE: Requires [Glide](https://glide.sh) along with [Go](https://golang.org/) version 1.6+.

## Usage

The package provides facilities to implement and connect live tile-based analytics to persistent in-memory storage services.

## Example

This minimalistic application shows how to register tile and meta data generators and connect them to a redis store.

```go
package main

import (
	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/generation/elastic"
	"github.com/unchartedsoftware/veldt/store/redis"
)

func main() {
	// Create pipeline
	pipeline := veldt.NewPipeline()

	// Add boolean expression types
	pipeline.Binary(elastic.NewBinaryExpression)
	pipeline.Unary(elastic.NewUnaryExpression)

	// Add query types to the pipeline
	pipeline.Query("exists", elastic.NewExists)
	pipeline.Query("has", elastic.NewHas)
	pipeline.Query("equals", elastic.NewEquals)
	pipeline.Query("range", elastic.NewRange)

	// Add tiles types to the pipeline
	pipeline.Tile("heatmap", elastic.NewHeatmapTile("localhost", "9200"))

	// Set the maximum concurrent tile requests
	pipeline.SetMaxConcurrent(32)
	// Set the tile requests queue length
	pipeline.SetQueueLength(1024)

	// Add a redis store to the pipeline
	pipeline.Store(redis.NewStore("localhost", "6379", -1))

	// register the pipeline
	veldt.Register("elastic", pipeline)

	// Create tile JSON request
	req :=
		`
		{
			"uri": "sample_index0",
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
					"right": 4294967296,
					"bottom": 0,
					"top": 4294967296,
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
						"gte": 18
					}
				}
			]
		}
		`

	// Generate a tile, this call will block until the tile is ready in the
	// store.
	err := veldt.GenerateTile("elastic", req)
	if err != nil {
		panic(err)
	}
	// Retrieve the tile form the store.
	tile, err := veldt.GetTileFromStore("elastic", t)
	if err != nil {
		panic(err)
	}
}
```

## Development

Clone the repository:

```bash
mkdir -p $GOPATH/src/github.com/unchartedsoftware
cd $GOPATH/src/github.com/unchartedsoftware
git clone git@github.com:unchartedsoftware/veldt.git
```

Install dependencies:

```bash
cd veldt
make install
```
