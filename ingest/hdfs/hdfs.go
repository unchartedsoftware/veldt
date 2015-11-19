package hdfs

import (
    "fmt"
    "os"

    "github.com/colinmarc/hdfs"
)

var hdfsClient *hdfs.Client

func getHdfsClient( host string, port string ) ( *hdfs.Client, error ) {
    endpoint := host + ":" + port
    if hdfsClient == nil {
        fmt.Println( "Connecting to HDFS: " + endpoint )
        client, err := hdfs.New( endpoint )
        hdfsClient = client
        return hdfsClient, err
    }
    return hdfsClient, nil
}

func Open( host string, port string, filePath string ) ( *hdfs.FileReader, error ) {
    // get hdfs client
    client, err := getHdfsClient( host, port )
    if err != nil {
        return nil, err
    }
    // get hdfs file reader
    return client.Open( filePath )
}

func ReadDir( host string, port string, path string ) ( []os.FileInfo, error ) {
    // get hdfs client
    client, err := getHdfsClient( host, port )
    if err != nil {
        return nil, err
    }
    fmt.Println( "Retreiving directory information from: " + path )
    return client.ReadDir( path )
}
