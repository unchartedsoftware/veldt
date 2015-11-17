package info

import (
    "fmt"
    "os"
    "runtime/debug"

    "github.com/unchartedsoftware/prism/ingest/hdfs"
)

type IngestInfo struct {
    Files []os.FileInfo
    NumTotalBytes int64
}

func GetIngestInfo( host string, port string, path string ) *IngestInfo {
    client, err := hdfs.GetHdfsClient( host, port )
    if err != nil {
        fmt.Println( err )
        debug.PrintStack()
        os.Exit(1)
    }
    fmt.Println( "Retreiving ingest directory information from: " + path )
    files, err := client.ReadDir( path )
    if err != nil {
        fmt.Println( err )
        debug.PrintStack()
        os.Exit(1)
    }
    var fileInfos []os.FileInfo
    var numTotalBytes int64 = 0
    for i:= 0; i<len( files ); i++ {
        file := files[i]
        if !file.IsDir() && file.Name() != ".SUCCESS" && file.Size() > 0 {
            // add to total bytes
            numTotalBytes += file.Size()
            // store file info
            fileInfos = append( fileInfos, file )
        }
    }
    return &IngestInfo{
        Files: fileInfos,
        NumTotalBytes: numTotalBytes,
    }
}
