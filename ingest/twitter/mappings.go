package twitter

func GetMappings() string {
    return `{
        "datum": {
            "properties": {
                "locality": {
                    "type": "object",
                    "properties": {
                        "location": {
                            "type": "geo_point"
                        },
                        "hashtags" : {
                          "type" : "string",
                          "index" : "not_analyzed"
                        },
                        "UserID" : {
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
