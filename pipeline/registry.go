package tile

import (
	"fmt"
	"runtime"
	"sync"
)

func (p *Pipeline) Query(id string, ctor prism.QueryCtor) {
	p.queries[id] = ctor
}

func (p *Pipeline) Tile(id string, ctor prism.TileCtor) {
	p.tiles[id] = ctor
}

func (p *Pipeline) Meta(id string, ctor prism.MetaCtor) {
	p.metas[id] = ctor
}

func (p *Pipeline) Store(ctor prism.StoreCtor) {
	p.store = ctor
}

func (p *Pipeline) GetQuery(id string, params map[string]interface{}) (prism.Query, error) {
	ctor, ok := p.queries[id]
	if !ok {
		return nil, fmt.Errorf("Unrecognized query type `%v`", id)
	}
	return ctor(params)
}

func (p *Pipeline) GetTile(id string, params map[string]interface{}) (prism.Tile, error) {
	ctor, ok := p.tiles[id]
	if !ok {
		return nil, fmt.Errorf("Unrecognized tile type `%v`", id)
	}
	return ctor(params)
}

func (p *Pipeline) GetMeta(id string, params map[string]interface{}) (prism.Meta, error) {
	ctor, ok := p.metas[id]
	if !ok {
		return nil, fmt.Errorf("Unrecognized meta type `%v`", id)
	}
	return ctor(params)
}

func (p *Pipeline) GetStore() {
	if p.store == nil {
		return nil, fmt.Errorf("No store has been provided")
	}
	return p.store(), nil
}
