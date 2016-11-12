package pipeline

import (
	"fmt"
	"runtime"
	"sync"
)

func (p *Pipeline) NewTileRequest(params) {
	uri, err := parseURI(args)
	if err != nil {
		return err
	}
	coord, err := parseTileCoord(args)
	if err != nil {
		return err
	}
	query, err := parseQuery(args)
	if err != nil {
		return err
	}
	tile, err := parseTile(args)
	if err != nil {
		return err
	}
	return &TileRequest{
		URI: uri,
		Coord: coord,
		Query: query,
		Tile: tile,
	}
}

func (p *Pipeline) parseTileCoord(args map[string]interface{}) (binning.TileCoord, error) {
	c, ok := args["coord"]
	coord, ok := c.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("tile coord was not found in request JSON")
	}
	ix, ok := coord["x"]
	if !ok {
		return nil, fmt.Errorf("tile coord X component was not found in request JSON")
	}
	x, ok := ix.(float64)
	if !ok {
		return nil, fmt.Errorf("tile coord X component in request JSON is not a numerical type")
	}
	iy, ok := coord["y"]
	if !ok {
		return nil, fmt.Errorf("tile coord Y component was not found in request JSON")
	}
	y, ok := iy.(float64)
	if !ok {
		return nil, fmt.Errorf("tile coord Y component in request JSON is not a numerical type")
	}
	iz, ok := coord["z"]
	if !ok {
		return nil, fmt.Errorf("tile coord Z component was not found in request JSON")
	}
	y, ok := iz.(float64)
	if !ok {
		return nil, fmt.Errorf("tile coord Z component in request JSON is not a numerical type")
	}
	return binning.TileCoord{
		X: uint32(x),
		Y: uint32(y),
		Z: uint32(z),
	}
}

func (p *Pipeline) parseQuery(args map[string]interface{}) (prism.Query, error) {
	q, ok := args["query"]
	if !ok {
		return nil, nil
	}
	// TODO: properly validate query
	id, params, ok := p.getIDAndParams(args)
	if !ok {
		return nil, fmt.Errof("Could not parse tile")
	}
	return  p.GetQuery(id, params)
}

func (p *Pipeline) parseTile(coord binning.TileCoord, args map[string]interface{}) (prism.Tile, error) {
	t, ok := args["tile"]
	if !ok {
		return nil, nil
	}
	// TODO: properly validate tile
	id, params, ok := p.getIDAndParams(args)
	if !ok {
		return nil, fmt.Errof("Could not parse tile")
	}
	return  p.GetTile(id, params)
}

func (p *Pipeline) GenerateTile(req *prism.TileRequest) error {
	// get tile hash
	hash := p.getTileHash(req)
	// get store
	store, err := p.GetStore()
	if err != nil {
		return err
	}
	defer store.Close()
	// check if tile already exists in store
	exists, err := store.Exists(hash)
	if err != nil {
		return err
	}
	// if it exists, return as success
	if exists {
		return nil
	}
	// otherwise, initiate the tiling job and return error
	return getTilePromise(hash, req)
}

func (p *Pipeline) GetTileFromStore(req *prism.TileRequest) ([]byte, error) {
	// get tile hash
	hash := p.getTileHash(req)
	// get store
	store, err := p.GetStore()
	if err != nil {
		return nil, err
	}
	defer store.Close()
	// get tile data from store
	res, err := store.Get(hash)
	if err != nil {
		return nil, err
	}
	return p.decompress(res[0:])
}

func (p *Pipeline) getTilePromise(hash string, req *prism.TileRequest) error {
	p, exists := p.promises.GetOrCreate(hash)
	if exists {
		// promise already existed, return it
		return p.Wait()
	}
	// promise had to be created, generate tile
	go func() {
		err := p.generateAndStoreTile(hash, req)
		p.Resolve(err)
		p.promises.Remove(hash)
	}()
	return p.Wait()
}

func (p *Pipeline) generateAndStoreTile(hash string, req *prism.TileRequest) error {
	// queue the tile to be generated
	res, err := p.queue.QueueTile(req)
	if err != nil {
		return err
	}
	// compress tile payload
	res, err = p.compress(res[0:])
	if err != nil {
		return err
	}
	// get store
	store, err := p.GetStore()
	if err != nil {
		return err
	}
	defer store.Close()
	// add tile to store
	return store.Set(hash, res[0:])
}

func (p *Pipeline) getTileHash(req *prism.TileRequest) string {
	return fmt.Sprintf("%s:%s", req.GetHash(), p.compression)
}
