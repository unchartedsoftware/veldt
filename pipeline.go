package veldt

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/unchartedsoftware/veldt/util/promise"
)

// Pipeline represents a cohesive tile and meta generation unit.
type Pipeline struct {
	queue       *queue
	queries     map[string]QueryCtor
	binary      QueryCtor
	unary       QueryCtor
	tiles       map[string]TileCtor
	metas       map[string]MetaCtor
	store       StoreCtor
	promises    *promise.Map
	compression string
}

// NewPipeline instantiates and returns a new pipeline struct.
func NewPipeline() *Pipeline {
	return &Pipeline{
		queue:       newQueue(),
		queries:     make(map[string]QueryCtor),
		tiles:       make(map[string]TileCtor),
		metas:       make(map[string]MetaCtor),
		promises:    promise.NewMap(),
		compression: "gzip",
	}
}

// SetMaxConcurrent sets the maximum concurrent tile requests allowed.
func (p *Pipeline) SetMaxConcurrent(max int) {
	p.queue.setMaxConcurrent(max)
}

// SetQueueLength sets the queue length for tiles to hold in the queue.
func (p *Pipeline) SetQueueLength(length int) {
	p.queue.setQueueLength(length)
}

// Query registers a query type under the provided ID string.
func (p *Pipeline) Query(id string, ctor QueryCtor) {
	p.queries[id] = ctor
}

// Binary registers a binary operator type under the provided ID string.
func (p *Pipeline) Binary(ctor QueryCtor) {
	p.binary = ctor
}

// Unary registers a unary operator type under the provided ID string.
func (p *Pipeline) Unary(ctor QueryCtor) {
	p.unary = ctor
}

// Tile registers a tile generation type under the provided ID string.
func (p *Pipeline) Tile(id string, ctor TileCtor) {
	p.tiles[id] = ctor
}

// Meta registers a metadata generation type under the provided ID string.
func (p *Pipeline) Meta(id string, ctor MetaCtor) {
	p.metas[id] = ctor
}

// Store registers the storage system used to cache generated data.
func (p *Pipeline) Store(ctor StoreCtor) {
	p.store = ctor
}

// GetQuery returns the instantiated query struct from the provided ID and JSON.
func (p *Pipeline) GetQuery(id string, args interface{}) (Query, error) {
	params, ok := args.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("`%s` is not of correct type", id)
	}
	ctor, ok := p.queries[id]
	if !ok {
		return nil, fmt.Errorf("unrecognized query type `%v`", id)
	}
	query, err := ctor()
	if err != nil {
		return nil, err
	}
	err = query.Parse(params)
	if err != nil {
		return nil, err
	}
	return query, nil
}

// GetBinary returns the instantiated binary operator struct from the provided
// ID and JSON.
func (p *Pipeline) GetBinary() (Query, error) {
	if p.binary == nil {
		return nil, fmt.Errorf("no binary query type has been provided")
	}
	return p.binary()
}

// GetUnary returns the instantiated unary operator struct from the provided
// ID and JSON.
func (p *Pipeline) GetUnary() (Query, error) {
	if p.unary == nil {
		return nil, fmt.Errorf("no unary query type has been provided")
	}
	return p.unary()
}

// GetTile returns the instantiated tile generator struct from the provided
// ID and JSON.
func (p *Pipeline) GetTile(id string, args interface{}) (Tile, error) {
	params, ok := args.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("`%s` is not of correct type", id)
	}
	ctor, ok := p.tiles[id]
	if !ok {
		return nil, fmt.Errorf("unrecognized tile type `%v`", id)
	}
	tile, err := ctor()
	if err != nil {
		return nil, err
	}
	err = tile.Parse(params)
	if err != nil {
		return nil, err
	}
	return tile, nil
}

// GetMeta returns the instantiated metedata generator struct from the provided
// ID and JSON.
func (p *Pipeline) GetMeta(id string, args interface{}) (Meta, error) {
	params, ok := args.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("`%s` is not of correct type", id)
	}
	ctor, ok := p.metas[id]
	if !ok {
		return nil, fmt.Errorf("unrecognized meta type `%v`", id)
	}
	meta, err := ctor()
	if err != nil {
		return nil, err
	}
	err = meta.Parse(params)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

// GetStore returns the instantiated store struct from the provided ID and JSON.
func (p *Pipeline) GetStore() (Store, error) {
	if p.store == nil {
		return nil, fmt.Errorf("no store type has been provided")
	}
	return p.store()
}

// NewTileRequest instantiates and returns a tile request struct from the
// provided JSON.
func (p *Pipeline) NewTileRequest(args map[string]interface{}) (*TileRequest, error) {
	// params are modified in place during validation, so create a copy
	copy, err := copyJSON(args)
	if err != nil {
		return nil, err
	}
	// validate request
	req, err := newValidator(p).validateTileRequest(copy)
	if err != nil {
		return nil, fmt.Errorf("invalid tile request:\n%s", err)
	}
	return req, nil
}

// GenerateTile generates a tile for the provided tile request.
func (p *Pipeline) GenerateTile(req *TileRequest) error {
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
	return p.getTilePromise(hash, req)
}

// GetTileFromStore retrieves the generated tile from the store.
func (p *Pipeline) GetTileFromStore(req *TileRequest) ([]byte, error) {
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
	return p.decompress(res)
}

func (p *Pipeline) getTilePromise(hash string, req *TileRequest) error {
	promise, exists := p.promises.GetOrCreate(hash)
	if exists {
		// promise already existed, return it
		return promise.Wait()
	}
	// promise had to be created, generate tile
	go func() {
		err := p.generateAndStoreTile(hash, req)
		promise.Resolve(err)
		p.promises.Remove(hash)
	}()
	return promise.Wait()
}

func (p *Pipeline) generateAndStoreTile(hash string, req *TileRequest) error {
	// queue the tile to be generated
	res, err := p.queue.queueTile(req)
	if err != nil {
		return err
	}
	// compress tile payload
	res, err = p.compress(res)
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
	return store.Set(hash, res)
}

func (p *Pipeline) getTileHash(req *TileRequest) string {
	return fmt.Sprintf("%s:%s", req.GetHash(), p.compression)
}

// NewMetaRequest instantiates and returns a metadata request struct from the
// provided JSON.
func (p *Pipeline) NewMetaRequest(args map[string]interface{}) (*MetaRequest, error) {
	// params are modified in place during validation, so create a copy
	copy, err := copyJSON(args)
	if err != nil {
		return nil, err
	}
	// validate request
	req, err := newValidator(p).validateMetaRequest(copy)
	if err != nil {
		return nil, fmt.Errorf("invalid meta request:\n%s", err)
	}
	return req, nil
}

// GenerateMeta generates metadata for the provided metadata request.
func (p *Pipeline) GenerateMeta(req *MetaRequest) error {
	// get tile hash
	hash := p.getMetaHash(req)
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
	return p.getMetaPromise(hash, req)
}

// GetMetaFromStore retrieves the generated tile from the store.
func (p *Pipeline) GetMetaFromStore(req *MetaRequest) ([]byte, error) {
	// get tile hash
	hash := p.getMetaHash(req)
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
	return p.decompress(res)
}

func (p *Pipeline) getMetaPromise(hash string, req *MetaRequest) error {
	promise, exists := p.promises.GetOrCreate(hash)
	if exists {
		// promise already existed, return it
		return promise.Wait()
	}
	// promise had to be created, generate tile
	go func() {
		err := p.generateAndStoreMeta(hash, req)
		promise.Resolve(err)
		p.promises.Remove(hash)
	}()
	return promise.Wait()
}

func (p *Pipeline) generateAndStoreMeta(hash string, req *MetaRequest) error {
	// queue the tile to be generated
	res, err := p.queue.queueMeta(req)
	if err != nil {
		return err
	}
	// compress tile payload
	res, err = p.compress(res)
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
	return store.Set(hash, res)
}

func (p *Pipeline) getMetaHash(req *MetaRequest) string {
	return fmt.Sprintf("%s:%s", req.GetHash(), p.compression)
}

func (p *Pipeline) compress(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	writer, ok := p.getWriter(&buffer)
	if ok {
		// compress
		_, err := writer.Write(data)
		if err != nil {
			return nil, err
		}
		err = writer.Close()
		if err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

func (p *Pipeline) decompress(data []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(data)
	reader, ok, err := p.getReader(buffer)
	if err != nil {
		return nil, err
	}
	if ok {
		// decompress
		var err error
		data, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		err = reader.Close()
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (p *Pipeline) getReader(buffer *bytes.Buffer) (io.ReadCloser, bool, error) {
	// use compression based reader if specified
	switch p.compression {
	case "gzip":
		reader, err := gzip.NewReader(buffer)
		return reader, true, err
	case "zlib":
		reader, err := zlib.NewReader(buffer)
		return reader, true, err
	default:
		return nil, false, nil
	}
}

func (p *Pipeline) getWriter(buffer *bytes.Buffer) (io.WriteCloser, bool) {
	// use compression based reader if specified
	switch p.compression {
	case "gzip":
		return gzip.NewWriter(buffer), true
	case "zlib":
		return zlib.NewWriter(buffer), true
	default:
		return nil, false
	}
}

func copyJSON(obj map[string]interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	var copy map[string]interface{}
	err = json.Unmarshal(bytes, &copy)
	if err != nil {
		return nil, err
	}
	return copy, nil
}
