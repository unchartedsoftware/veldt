package info

import (
    "fmt"
    "os"
    "runtime/debug"

    "github.com/unchartedsoftware/prism/ingest/hdfs"
)

type IngestInfo struct {
    Files []os.FileInfo
    NumTotalBytes uint64
}

func GetIngestInfo( host string, port string, path string ) *IngestInfo {
    files, err := hdfs.ReadDir( host, port, path )
    if err != nil {
        fmt.Println( err )
        debug.PrintStack()
        os.Exit(1)
    }
    var fileInfos []os.FileInfo
    var numTotalBytes int64 = 0
    for _, file := range files {
        if !file.IsDir() && file.Name() != ".SUCCESS" && file.Size() > 0 {
            // add to total bytes
            numTotalBytes += file.Size()
            // store file info
            fileInfos = append( fileInfos, file )
        }
    }
    return &IngestInfo{
        Files: fileInfos,
        NumTotalBytes: uint64( numTotalBytes ),
    }
}
