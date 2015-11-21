package info

import (
	"os"

	"github.com/unchartedsoftware/prism/ingest/hdfs"
)

// IngestInfo represents a directory worth of ingestible data.
type IngestInfo struct {
	Files         []os.FileInfo
	NumTotalBytes uint64
}

// GetIngestInfo returns an array of os.FileInfo and the total number of bytes in the provided directory.
func GetIngestInfo(host string, port string, path string) (*IngestInfo, error) {
	files, err := hdfs.ReadDir(host, port, path)
	if err != nil {
		return nil, err
	}
	var fileInfos []os.FileInfo
	var numTotalBytes int64
	for _, file := range files {
		if !file.IsDir() && file.Name() != ".SUCCESS" && file.Size() > 0 {
			// add to total bytes
			numTotalBytes += file.Size()
			// store file info
			fileInfos = append(fileInfos, file)
		}
	}
	return &IngestInfo{
		Files:         fileInfos,
		NumTotalBytes: uint64(numTotalBytes),
	}, nil
}
