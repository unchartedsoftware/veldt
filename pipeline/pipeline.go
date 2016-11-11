package tile

import (
	"fmt"
	"runtime"
	"sync"
)

type Pipeline struct {
	queue *Queue
	queries map[string]prism.QueryCtor
	tiles map[string]prism.TileCtor
	metas map[string]prism.MetaCtor
	store prism.StoreCtor
	promises promise.Map
	compression string
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		queue: NewQueue(),
		queries: make(map[string]prism.QueryCtor),
		tiles: make(map[string]prism.TileCtor),
		metas: make(map[string]prism.MetaCtor),
		promises: promise.NewMap(),
		compression: "gzip",
	}
}

// SetMaxConcurrent sets the maximum concurrent tile requests allowed.
func (p *Pipeline) SetMaxConcurrent(max int) {
	p.Queue.SetMaxConcurrent(max)
}

// SetQueueLength sets the queue length for tiles to hold in the queue.
func (p *Pipeline) SetQueueLength(length int) {
	p.Queue.SetQueueLength(length)
}

func (p *Pipeline) getIDAndParams(args map[string]interface{}) (string, map[string]interface{}, bool) {
	var key string
	var value map[string]interface{}
	found := false
	for k, v := range args {
		val, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		key = k
		value = val
		found = true
		break
	}
	return key, value, found
}
