package info

import (
	"fmt"
	"os"

	"github.com/unchartedsoftware/prism/ingest/hdfs"
	"github.com/unchartedsoftware/prism/util"
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
	return file.IsDir() &&
		file.Name() != "es-mapping-files"
}

func isValidFile(file os.FileInfo) bool {
	return !file.IsDir() &&
		file.Name() != ".SUCCESS" &&
		file.Name() != "_SUCCESS" &&
		file.Size() > 0
}

// GetIngestInfo returns an array of os.FileInfo and the total number of bytes in the provided directory.
func GetIngestInfo(host string, port string, path string) (*IngestInfo, error) {
	// create empty slice of fileinfos to populate
	files, err := hdfs.ReadDir(host, port, path)
	if err != nil {
		return nil, err
	}
	// data to populate
	var ingestFiles []IngestFile
	numTotalBytes := uint64(0)
	// for each file / dir
	for _, file := range files {
		if isValidDir(file) {
			// depth-first traversal into sub directories
			subInfo, err := GetIngestInfo(host, port, path+"/"+file.Name())
			ingestFiles = append(ingestFiles, subInfo.Files...)
			numTotalBytes += subInfo.NumTotalBytes
			// if any errors, append them
			if err != nil {
				return nil, err
			}
			fmt.Printf("Retreiving ingest information from: %s, %d files containing %s\n",
				path+"/"+file.Name(),
				len(subInfo.Files),
				util.FormatBytes(float64(subInfo.NumTotalBytes)))
		} else if isValidFile(file) {
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
	return &IngestInfo{
		Files:         ingestFiles[0:],
		NumTotalBytes: numTotalBytes,
	}, nil
}
