package main

import (
	"runtime"
	"syscall"

	"github.com/zenazn/goji/graceful"

	"github.com/unchartedsoftware/prism/server/api"
	"github.com/unchartedsoftware/prism/server/conf"
	"github.com/unchartedsoftware/prism/tiling"
	"github.com/unchartedsoftware/prism/tiling/elastic"
	"github.com/unchartedsoftware/prism/util/log"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	// parse flags into config struct
	config := conf.ParseCommandLine()
	log.Debugf("Prism serving from directory '%s'", config.Public)
	log.Debugf("Prism server listening on port %s", config.Port)

	// register available tiling types
	tiling.Register("heatmap", elastic.GetHeatmapTile, elastic.GetMeta)
	tiling.Register("topiccount", elastic.GetTopicCountTile, elastic.GetMeta)
	tiling.Register("crossplot_deprecated", elastic.GetCrossPlotDeprecatedTile, elastic.GetMeta)

	meta, err := elastic.GetMeta("http://10.64.16.120:9200", "isil_twitter_weekly")
	if err != nil {
		log.Error(err)
	}
	for k, v := range meta {
		if v.Extrema != nil {
			log.Debugf("%s: type: %s, min: %f, max: %f", k, v.Type, v.Extrema.Min, v.Extrema.Max)
		} else {
			log.Debugf("%s: type: %s", k, v.Type)
		}
	}

	// start the web server
	graceful.AddSignal(syscall.SIGINT, syscall.SIGTERM)

	// create server
	app := api.New()
	err = graceful.ListenAndServe(":"+config.Port, app)
	if err != nil {
		log.Error(err)
	}
	graceful.Wait()
}
