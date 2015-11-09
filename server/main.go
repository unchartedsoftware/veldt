package main

import (
    "flag"
	"log"
	"runtime"
	"syscall"

	"github.com/zenazn/goji/graceful"

    "github.com/unchartedsoftware/prism/server/conf"
    "github.com/unchartedsoftware/prism/server/api"
)

var (
	port = flag.CommandLine.String( "port", "8080", "Port to bind HTTP server" )
	prod = flag.CommandLine.Bool( "prod", false, "Production flag" )
)

func main() {

    runtime.GOMAXPROCS( runtime.NumCPU() )

    // Parse the flags and store them as a conf struct
	configuration := conf.Conf{
		Prod: *prod,
        Port: *port,
	}
	conf.SaveConf( &configuration )

    // Start the web server
	graceful.AddSignal( syscall.SIGINT, syscall.SIGTERM )
	app := api.New()
	log.Println( "Prism server listening on port " + *port )
	err := graceful.ListenAndServe( ":" + *port, app )
	if err != nil {
		log.Fatal( err )
	}
	graceful.Wait()
}
