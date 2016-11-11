package gzip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/unchartedsoftware/prism/store"
)

func (p *Pipeline) compress(compression string, data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	writer, err := getWriter(compression, &buffer)
	if err != nil {
		return nil, err
	}
	_, err = writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes()[0:], nil
}

func (p *Pipeline) decompress(typ string, data []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(data[0:])
	reader, err := getReader(compression, &buffer)
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

func getReader(compression string, buffer *bytes.Buffer) (io.Reader, error) {
	// use compression based reader if specified
	switch compression {
	case "gzip":
		return gzip.NewReader(buffer)
	case "bzip2":
		return bzip2.NewReader(buffer), nil
	case "flate":
		return flate.NewReader(buffer), nil
	case "zlib":
		return zlib.NewReader(buffer)
	default:
		return buffer, nil
	}
}

func getWriter(compression string, buffer *bytes.Buffer) (io.Reader, error) {
	// use compression based reader if specified
	switch compression {
	case "gzip":
		return gzip.NewWriter(buffer)
	case "bzip2":
		return bzip2.NewWriter(buffer), nil
	case "flate":
		return flate.NewWriter(buffer), nil
	case "zlib":
		return zlib.NewWriter(buffer)
	default:
		return buffer, nil
	}
}
