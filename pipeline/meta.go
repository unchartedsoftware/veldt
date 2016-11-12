package pipeline

import (
	"fmt"
	"runtime"
	"sync"
)

func (p *Pipeline) NewMetaRequest(params) {
	uri, err := parseURI(args)
	if err != nil {
		return err
	}
	meta, err := parseMeta(args)
	if err != nil {
		return err
	}
	return &MetaRequest{
		URI: uri,
		Meta: meta,
	}
}

func (p *Pipeline) parseMeta(coord binning.MetaCoord, args map[string]interface{}) (prism.Meta, error) {
	t, ok := args["meta"]
	if !ok {
		return nil, nil
	}
	// TODO: properly validate tile
	id, params, ok := p.getIDAndParams(args)
	if !ok {
		return nil, fmt.Errof("Could not parse meta")
	}
	return  p.GetMeta(id, params)
}

func (p *Pipeline) GenerateMeta(req *prism.MetaRequest) error {
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
	return getMetaPromise(hash, req)
}

func (p *Pipeline) GetMetaFromStore(req *prism.MetaRequest) ([]byte, error) {
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

func (p *Pipeline) getMetaPromise(hash string, req *prism.MetaRequest) error {
	p, exists := p.promises.GetOrCreate(hash)
	if exists {
		// promise already existed, return it
		return p.Wait()
	}
	// promise had to be created, generate tile
	go func() {
		err := p.generateAndStoreMeta(hash, req)
		p.Resolve(err)
		p.promises.Remove(hash)
	}()
	return p.Wait()
}

func (p *Pipeline) generateAndStoreMeta(hash string, req *prism.MetaRequest, store prism.Store) error {
	// queue the tile to be generated
	res, err := p.queue.QueueMeta(req)
	if err != nil {
		return err
	}
	// compress tile payload
	res, err = p.compress(p.compression, res[0:])
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

func (p *Pipeline) getMetaHash(req *prism.MetaRequest) string {
	return fmt.Sprintf("%s:%s", req.GetHash(), p.compression)
}
