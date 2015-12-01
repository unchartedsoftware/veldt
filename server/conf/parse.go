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
	public := flag.String("publicDir", "./build", "The public directory to static serve from")

	// Parse the flags
	flag.Parse()

	// Set and save config
	config := &Conf{
		Port:       *port,
		Public:     *public,
		Aliases:    aliases,
		InvAliases: invertAliases(aliases),
	}
	SaveConf(config)
	return config
}
