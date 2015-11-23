package main

import (
	"flag"
	"runtime"
	"syscall"

	"github.com/zenazn/goji/graceful"

	"github.com/unchartedsoftware/prism/server/api"
	"github.com/unchartedsoftware/prism/server/conf"
	"github.com/unchartedsoftware/prism/util/log"
)

var (
	port   = flag.CommandLine.String("port", "8080", "Port to bind HTTP server")
	prod   = flag.CommandLine.Bool("prod", false, "Production flag")
	public = flag.CommandLine.String("publicDir", "./build", "The public directory to static serve from")
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	// Parse the flags and store them as a conf struct
	flag.Parse()
	config := conf.Conf{
		Prod:   *prod,
		Port:   *port,
		Public: *public,
	}
	conf.SaveConf(&config)

	// Start the web server
	graceful.AddSignal(syscall.SIGINT, syscall.SIGTERM)
	app := api.New()
	log.Debug("Prism server listening on port " + config.Port)
	err := graceful.ListenAndServe(":"+config.Port, app)
	if err != nil {
		log.Error(err)
	}
	graceful.Wait()
}
