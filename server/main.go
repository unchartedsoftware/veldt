package main

import (
	"runtime"
	"syscall"

	"github.com/zenazn/goji/graceful"

	"github.com/unchartedsoftware/prism/server/api"
	"github.com/unchartedsoftware/prism/server/conf"
	"github.com/unchartedsoftware/prism/util/log"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	// parse flags into config struct
	config := conf.ParseCommandLine()

	log.Debugf("Prism serving from directory '%s'", config.Public)
	log.Debugf("Prism server listening on port %s", config.Port)

	// start the web server
	graceful.AddSignal(syscall.SIGINT, syscall.SIGTERM)

	// create server
	app := api.New()
	err := graceful.ListenAndServe(":"+config.Port, app)
	if err != nil {
		log.Error(err)
	}
	graceful.Wait()
}
