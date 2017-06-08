package veldt

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/unchartedsoftware/veldt/util/json"
	"github.com/unchartedsoftware/veldt/util/promise"
	"github.com/unchartedsoftware/veldt/util/queue"
)

// Pipeline represents a cohesive tile and meta generation unit.
type Pipeline struct {
	queue       *queue.Queue
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
		queue:       queue.NewQueue(),
		queries:     make(map[string]QueryCtor),
		tiles:       make(map[string]TileCtor),
		metas:       make(map[string]MetaCtor),
		promises:    promise.NewMap(),
		compression: "gzip",
	}
}

// SetMaxConcurrent sets the maximum concurrent tile requests allowed.
func (p *Pipeline) SetMaxConcurrent(max int) {
	p.queue.SetMaxConcurrent(max)
}

// SetQueueLength sets the queue length for tiles to hold in the queue.
func (p *Pipeline) SetQueueLength(length int) {
	p.queue.SetLength(length)
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

// GetHash returns a unique hash for the state of the pipeline.
func (p *Pipeline) GetHash() string {
	return p.compression
}

// NewTileRequest instantiates and returns a tile request struct from the
// provided JSON.
func (p *Pipeline) NewTileRequest(args map[string]interface{}) (*TileRequest, error) {
	// params are modified in place during validation, so create a copy
	copy, err := json.Copy(args)
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

// NewMetaRequest instantiates and returns a metadata request struct from the
// provided JSON.
func (p *Pipeline) NewMetaRequest(args map[string]interface{}) (*MetaRequest, error) {
	// params are modified in place during validation, so create a copy
	copy, err := json.Copy(args)
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

// Generate generates data for the provided request.
func (p *Pipeline) Generate(req Request) error {
	// get hash
	hash := p.getHash(req)
	// get store
	store, err := p.GetStore()
	if err != nil {
		return err
	}
	defer store.Close()
	// check if already exists in store
	exists, err := store.Exists(hash)
	if err != nil {
		return err
	}
	// if it exists, return as success
	if exists {
		return nil
	}
	// otherwise, initiate the generation task and return error
	return p.getPromise(hash, req)
}

// Get retrieves the generated data from the store.
func (p *Pipeline) Get(req Request) ([]byte, error) {
	// get hash
	hash := p.getHash(req)
	// get store
	store, err := p.GetStore()
	if err != nil {
		return nil, err
	}
	defer store.Close()
	// get data from store
	res, err := store.Get(hash)
	if err != nil {
		return nil, err
	}
	return p.decompress(res)
}

// GenerateAndGet retrieves the generated data from the store, if it
// does not exist, generate it before retrieval.
func (p *Pipeline) GenerateAndGet(req Request) ([]byte, error) {
	// get hash
	hash := p.getHash(req)
	// get store
	store, err := p.GetStore()
	if err != nil {
		return nil, err
	}
	defer store.Close()
	// check if already exists in store
	exists, err := store.Exists(hash)
	if err != nil {
		return nil, err
	}
	// check if it exists
	if !exists {
		// if not, initiate the tiling job
		err = p.getPromise(hash, req)
		if err != nil {
			return nil, err
		}
	}
	// get data from store
	res, err := store.Get(hash)
	if err != nil {
		return nil, err
	}
	return p.decompress(res)
}

func (p *Pipeline) getPromise(hash string, req Request) error {
	promise, exists := p.promises.GetOrCreate(hash)
	if exists {
		// promise already existed, return it
		return promise.Wait()
	}
	// promise had to be created, generate data
	go func() {
		err := p.generateAndStore(hash, req)
		promise.Resolve(err)
		p.promises.Remove(hash)
	}()
	return promise.Wait()
}

func (p *Pipeline) generateAndStore(hash string, req Request) error {
	// queue the tile to be generated
	res, err := p.queue.Send(req)
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

func (p *Pipeline) getHash(req Request) string {
	return fmt.Sprintf("%s:%s", req.GetHash(), p.GetHash())
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

func (p *Pipeline) getReader(buffer io.Reader) (io.ReadCloser, bool, error) {
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

func (p *Pipeline) getWriter(buffer io.Writer) (io.WriteCloser, bool) {
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
