package main

import (
	"runtime"
	"syscall"

	"github.com/zenazn/goji/graceful"

	"github.com/unchartedsoftware/prism/generation/elastic"
	"github.com/unchartedsoftware/prism/generation/meta"
	"github.com/unchartedsoftware/prism/generation/tile"
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

	// register available tiling types
	tile.Register("heatmap", elastic.NewHeatmapTile)
	tile.Register("topiccount", elastic.NewTopicCountTile)
	tile.Register("topicfrequency", elastic.NewTopicFrequencyTile)

	// register available meta data types
	meta.Register("default", elastic.NewDefaultMeta)

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
