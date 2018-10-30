<img width="600" src="https://rawgit.com/unchartedsoftware/veldt/master/logo.svg" alt="veldt" />

> Scalable on-demand tile-based analytics

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](http://godoc.org/github.com/unchartedsoftware/veldt)
[![Build Status](https://travis-ci.org/unchartedsoftware/veldt.svg?branch=master)](https://travis-ci.org/unchartedsoftware/veldt)
[![Go Report Card](https://goreportcard.com/badge/github.com/unchartedsoftware/veldt)](https://goreportcard.com/report/github.com/unchartedsoftware/veldt)

## Dependencies

Requires the [Go](https://golang.org/) programming language binaries with the `GOPATH` environment variable specified and `$GOPATH/bin` in your `PATH`.

## Installation

NOTE: Requires [Go](https://golang.org/) version 1.11+.

```bash
go get github.com/unchartedsoftware/veldt
```

## Usage

The package provides facilities to implement and connect live tile-based analytics to persistent in-memory storage services.

## Example

This minimalistic application shows how to register tile and meta data generators and connect them to a redis store.

```go
package main

import (
	"encoding/json"

	"github.com/unchartedsoftware/plog"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/generation/elastic"
	"github.com/unchartedsoftware/veldt/store/redis"
)

func JSON(str string) map[string]interface{} {
	var j map[string]interface{}
	err := json.Unmarshal([]byte(str), &j)
	if err != nil {
		panic(err)
	}
	return j
}

func main() {

	// Set logger for internal warnings and errors
	logger := log.NewLogger()
	veldt.SetLogger(veldt.Warn, logger)

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

	// Create tile JSON request
	arg := JSON(
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
		`)

	// Generate a tile directly from the pipeline

	// Instantiate a request object
	req, err := pipeline.NewTileRequest(arg)
	if err != nil {
		panic(err)
	}

	// Generate the tile, this call will block until the tile is ready in the
	// store.
	err = pipeline.Generate(req)
	if err != nil {
		panic(err)
	}

	// Retrieve the tile from the store
	bytes, err := pipeline.GetFromStore(req)
	if err != nil {
		panic(err)
	}

	// Generate a tile by pipeline ID (useful when using HTTP / WebSocket API)

	// Register the pipeline under a string ID
	veldt.Register("elastic", pipeline)

	// Generate a tile, this call will block until the tile is ready in the
	// store.
	err = veldt.GenerateTile("elastic", arg)
	if err != nil {
		panic(err)
	}

	// Retrieve the tile from the store
	bytes, err = veldt.GetTileFromStore("elastic", arg)
	if err != nil {
		panic(err)
	}
}
```

## Development

NOTE: Requires [Go](https://golang.org/) version 1.11+.

Clone the repository outside of your `$GOPATH`:

```bash
git clone git@github.com:unchartedsoftware/veldt.git
cd veldt
```

Install dependencies:

```bash
make install
```
