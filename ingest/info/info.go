package info

import (
	"os"

	"github.com/unchartedsoftware/prism/ingest/es"
	"github.com/unchartedsoftware/prism/ingest/hdfs"
	"github.com/unchartedsoftware/prism/util"
	"github.com/unchartedsoftware/prism/log"
)

// IngestInfo represents a directory worth of ingestible data.
type IngestInfo struct {
	Files         []IngestFile
	NumTotalBytes uint64
}

// IngestFile represents a single file to ingest.
type IngestFile struct {
	Name string
	Path string
	Size uint64
}

func isValidDir(file os.FileInfo) bool {
	return file.IsDir()
}

func isValidFile(file os.FileInfo) bool {
	return !file.IsDir() &&
		file.Name() != "_SUCCESS" &&
		file.Size() > 0
}

// GetIngestInfo returns an array of os.FileInfo and the total number of bytes in the provided directory.
func GetIngestInfo(host string, port string, path string, document es.Document) (*IngestInfo, error) {
	// read directory
	files, err := hdfs.ReadDir(host, port, path)
	if err != nil {
		return nil, err
	}
	// data to populate
	var ingestFiles []IngestFile
	numTotalBytes := uint64(0)
	log.Debugf("Retreiving ingest info from: %s", path)
	// for each file / dir
	for _, file := range files {
		if isValidDir(file) && document.FilterDir(file.Name()) {
			// depth-first traversal into sub directories
			subInfo, err := GetIngestInfo(host, port, path+"/"+file.Name(), document)
			if err != nil {
				return nil, err
			}
			ingestFiles = append(ingestFiles, subInfo.Files...)
			numTotalBytes += subInfo.NumTotalBytes
		} else if isValidFile(file) && document.FilterFile(file.Name()) {
			// add to total bytes
			numTotalBytes += uint64(file.Size())
			// store file info
			ingestFiles = append(ingestFiles, IngestFile{
				Name: file.Name(),
				Path: path,
				Size: uint64(file.Size()),
			})
		}
	}
	// print if we have found files
	if len(ingestFiles) > 0 {
		log.Debugf("Found %d files containing %s of ingestible data",
			len(ingestFiles),
			util.FormatBytes(float64(numTotalBytes)))
	}
	// return ingest info
	return &IngestInfo{
		Files:         ingestFiles[0:],
		NumTotalBytes: numTotalBytes,
	}, nil
}
