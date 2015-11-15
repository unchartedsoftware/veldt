package hdfs

import (
    "fmt"

    "github.com/colinmarc/hdfs"
)

var hdfsClient *hdfs.Client

func GetHdfsClient( host string, port string ) ( *hdfs.Client, error ) {
    endpoint := host + ":" + port
    if hdfsClient == nil {
        fmt.Println( "Connecting to HDFS: " + endpoint )
        client, err := hdfs.New( endpoint )
        hdfsClient = client
        return hdfsClient, err
    }
    return hdfsClient, nil
}
