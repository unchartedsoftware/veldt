package conf

import (
	"flag"
)

// ParseCommandLine parses the commandline arguments and returns a Conf object.
func ParseCommandLine() *Conf {
	// get alias flags
	aliases := make(AliasMap)
	flag.Var(&aliases, "alias", "Alias a string value by another string.")

	port := flag.String("port", "8080", "Port to bind HTTP server")
	public := flag.String("public", "./build/public", "The public directory to static serve from")

	redisHost := flag.String("redis-host", "localhost", "Host to connect to redis server")
	redisPort := flag.String("redis-port", "6379", "Port to connect to Redis server")

	// Parse the flags
	flag.Parse()

	// Set and save config
	config := &Conf{
		Port:       *port,
		Public:     *public,
		Aliases:    aliases,
		InvAliases: invertAliases(aliases),
		RedisHost:  *redisHost,
		RedisPort:  *redisPort,
	}
	SaveConf(config)
	return config
}
