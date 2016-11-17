package prism

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/unchartedsoftware/prism/util/promise"
)

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

func (p *Pipeline) Query(id string, ctor QueryCtor) {
	p.queries[id] = ctor
}

func (p *Pipeline) Binary(ctor QueryCtor) {
	p.binary = ctor
}

func (p *Pipeline) Unary(ctor QueryCtor) {
	p.unary = ctor
}

func (p *Pipeline) Tile(id string, ctor TileCtor) {
	p.tiles[id] = ctor
}

func (p *Pipeline) Meta(id string, ctor MetaCtor) {
	p.metas[id] = ctor
}

func (p *Pipeline) Store(ctor StoreCtor) {
	p.store = ctor
}

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

func (p *Pipeline) GetBinary() (Query, error) {
	if p.binary == nil {
		return nil, fmt.Errorf("no binary query type has been provided")
	}
	return p.binary()
}

func (p *Pipeline) GetUnary() (Query, error) {
	if p.unary == nil {
		return nil, fmt.Errorf("no unary query type has been provided")
	}
	return p.unary()
}

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

func (p *Pipeline) GetStore() (Store, error) {
	if p.store == nil {
		return nil, fmt.Errorf("no store type has been provided")
	}
	return p.store()
}

func (p *Pipeline) NewTileRequest(args map[string]interface{}) (*TileRequest, error) {
	validator := NewValidator(p)
	req, err := validator.ValidateTileRequest(args)
	if err != nil {
		return nil, fmt.Errorf("invalid tile request:\n%s", err)
	}
	return req, nil
}

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
	return p.decompress(res[0:])
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

func (p *Pipeline) getTileHash(req *TileRequest) string {
	return fmt.Sprintf("%s:%s", req.GetHash(), p.compression)
}

func (p *Pipeline) NewMetaRequest(args map[string]interface{}) (*MetaRequest, error) {
	validator := NewValidator(p)
	req, err := validator.ValidateMetaRequest(args)
	if err != nil {
		return nil, fmt.Errorf("invalid meta request:\n%s", err)
	}
	return req, nil
}

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
	return p.decompress(res[0:])
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
	return buffer.Bytes()[0:], nil
}

func (p *Pipeline) decompress(data []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(data[0:])
	reader, err, ok := p.getReader(buffer)
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
	return data[0:], nil
}

func (p *Pipeline) getReader(buffer *bytes.Buffer) (io.ReadCloser, error, bool) {
	// use compression based reader if specified
	switch p.compression {
	case "gzip":
		reader, err := gzip.NewReader(buffer)
		return reader, err, true
	case "zlib":
		reader, err := zlib.NewReader(buffer)
		return reader, err, true
	default:
		return nil, nil, false
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
