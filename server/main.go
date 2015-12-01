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

	// start the web server
	graceful.AddSignal(syscall.SIGINT, syscall.SIGTERM)
	app := api.New()
	log.Debug("Prism server listening on port " + config.Port)
	err := graceful.ListenAndServe(":"+config.Port, app)
	if err != nil {
		log.Error(err)
	}
	graceful.Wait()
}
