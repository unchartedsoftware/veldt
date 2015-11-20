package twitter

import (
    "github.com/unchartedsoftware/prism/ingest/conf"
)

// GetMappings returns the mappings object for twitter elasticsearch data.
func GetMappings() string {
    config := conf.GetConf()
    return `{
        "` + config.EsType + `": {
            "properties": {
                "locality": {
                    "type": "object",
                    "properties": {
                        "location": {
                            "type": "geo_point"
                        },
                        "userid" : {
                          "type" : "string",
                          "index" : "not_analyzed"
                        },
                        "username" : {
                          "type" : "string",
                          "index" : "not_analyzed"
                        }
                    }
                }
            }
        }
    }`
}
