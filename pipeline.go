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

func (p *Pipeline) Tile(id string, ctor TileCtor) {
	p.tiles[id] = ctor
}

func (p *Pipeline) Meta(id string, ctor MetaCtor) {
	p.metas[id] = ctor
}

func (p *Pipeline) Store(ctor StoreCtor) {
	p.store = ctor
}

func (p *Pipeline) GetQuery(id string, params map[string]interface{}) (Query, error) {
	ctor, ok := p.queries[id]
	if !ok {
		return nil, fmt.Errorf("Unrecognized query type `%v`", id)
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

func (p *Pipeline) GetTile(id string, params map[string]interface{}) (Tile, error) {
	ctor, ok := p.tiles[id]
	if !ok {
		return nil, fmt.Errorf("Unrecognized tile type `%v`", id)
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

func (p *Pipeline) GetMeta(id string, params map[string]interface{}) (Meta, error) {
	ctor, ok := p.metas[id]
	if !ok {
		return nil, fmt.Errorf("Unrecognized meta type `%v`", id)
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
		return nil, fmt.Errorf("No store has been provided")
	}
	return p.store()
}

func (p *Pipeline) NewTileRequest(args map[string]interface{}) (*TileRequest, error) {
	validator := NewValidator(p)
	return validator.ValidateTileRequest(args)
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
	return validator.ValidateMetaRequest(args)
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
	writer, err := p.getWriter(&buffer)
	if err != nil {
		return nil, err
	}
	_, err = writer.Write(data)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes()[0:], nil
}

func (p *Pipeline) decompress(data []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(data[0:])
	reader, err := p.getReader(buffer)
	if err != nil {
		return nil, err
	}
	data, err = ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data[0:], nil
}

func (p *Pipeline) getReader(buffer *bytes.Buffer) (io.Reader, error) {
	// use compression based reader if specified
	switch p.compression {
	case "gzip":
		return gzip.NewReader(buffer)
	case "zlib":
		return zlib.NewReader(buffer)
	default:
		return buffer, nil
	}
}

func (p *Pipeline) getWriter(buffer *bytes.Buffer) (io.Writer, error) {
	// use compression based reader if specified
	switch p.compression {
	case "gzip":
		return gzip.NewWriter(buffer), nil
	case "zlib":
		return zlib.NewWriter(buffer), nil
	default:
		return buffer, nil
	}
}
