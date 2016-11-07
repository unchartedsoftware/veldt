package gzip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/unchartedsoftware/prism/store"
)

// Compressor represents a gzip compressor.
type Compressor struct {
}

// NewCompressor instantiates and returns a new gzip compressor instance.
func NewCompressor() store.CompressorConstructor {
	return func() store.Compressor {
		return &Compressor{}
	}
}

// Compress compresses the provided bytes.
func (c *Compressor) Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes()[0:], nil
}

// Decompress decompresses the provided bytes.
func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	b := bytes.NewBuffer(data[0:])
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	data, err = ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	err = r.Close()
	if err != nil {
		return nil, err
	}
	return data[0:], nil
}
