package salt

import (
	"fmt"
	"strings"
	"io/ioutil"
	"github.com/liyinhgqw/typesafe-config/parse"
)

// QueueConfiguration describes how a specific queue is to be created
type QueueConfiguration struct {
	queue     string
	durable   bool
	deletable bool
	exclusive bool
	noWait    bool
}

// Configuration describes how the salt connection is to be created
type Configuration struct {
	host                string
	port                int64
	serverQueue         string
	queueConfigurations map[string]*QueueConfiguration
}

func getConfigString (key string, config *parse.Config) Option {
	result, err := config.GetString(key)
	if err != nil {
		return Option{false, nil}
	}
	return Option{true, result}
}
func getConfigInt (key string, config *parse.Config) Option {
	result, err := config.GetInt(key)
	if err != nil {
		return Option{false, nil}
	}
	return Option{true, result}
}
func getConfigBool (key string, config *parse.Config) Option {
	result, err := config.GetBool(key)
	if err != nil {
		return Option{false, nil}
	}
	return Option{true, result}
}
func getConfigArray (key string, config *parse.Config) Option {
	result, err := config.GetBool(key)
	if err != nil {
		return Option{false, nil}
	}
	return Option{true, result}
}
func getSubConfig (key string, config *parse.Config) Option {
	result, err := config.GetValue(key)
	if err != nil {
		return Option{false, nil}
	}
	return Option{true, result}
}

func stripTerminalQuotes (text string) string {
	// TrimPrefix and TrimSuffix already work conditionally - we have our outer
	// condition here because we only want to remove the prefix and suffix if
	// both are present at the same time.
	if strings.HasPrefix(text, "\"") && strings.HasSuffix(text, "\"") {
		return strings.TrimPrefix(strings.TrimSuffix(text, "\""), "\"")
	}
	return text
}

// getKeys gets a list of all keys found under the given node.  This really should get
// pushed back to https://github.com/liyinhgqw/typesafe-config
func getKeys (root parse.Node, path ...string) ([]string, error) {
	if root.Type() != parse.NodeMap {
		return nil, fmt.Errorf("Root node is not a map")
	}
	mapRoot := root.(*parse.MapNode)

	// No path - return all nodes
	if 0 == len(path) {
		keys := make([]string, len(mapRoot.Nodes))
		i := 0
		for k := range mapRoot.Nodes {
			keys[i] = k
			i++
		}
		return keys, nil
	}

	var elem, contains = mapRoot.Nodes[path[0]]

	if !contains {
		return nil, fmt.Errorf("Node %s not found in map", path[0])
	}
	return getKeys(elem, path[1:]...)
}

// ReadConfiguration reads a typesafe-config configuration file that sets up our salt connection
func ReadConfiguration (filename string) (*Configuration, error) {
	// Read the config file
	configString, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// Parse it out into a config map
	configMap, err := parse.Parse("salt", string(configString))
	if err != nil {
		panic(err)
	}

	// We need to get the list of queues from the map, before it is made into a config
	queueKeys, err := getKeys(configMap.Root, "rabbitmq", "queues")
	if err != nil {
		// No keys found; initialize an empty list of pre-initialized queues
		queueKeys = make([]string, 0)
	}

	// Convert the map to a config object for further processing
	config := configMap.GetConfig()
	// Get the host from the config
	host := stripTerminalQuotes(getConfigString("rabbitmq.host", config).OrElse("localhost").(string))
	// Get the port from the config
	port := getConfigInt("rabbitmq.port", config).OrElse(int64(5672)).(int64)
	// Get the server queue from the config
	serverQueue := stripTerminalQuotes(getConfigString("communications.queue", config).OrElse("salt").(string))
	// Get any queues that need to be pre-initialized
	queues := make(map[string]*QueueConfiguration)
	for _, queueKey := range queueKeys {
		queueConfig := getSubConfig("rabbitmq.queues."+queueKey, config).OrElse(&parse.Config{}).(*parse.Config)
		queueInterface, err := getConfigString("name", queueConfig).Value()
		if nil != err {
			panic(err)
		}
		queue := stripTerminalQuotes(queueInterface.(string))
		durable := getConfigBool("durable", queueConfig).OrElse(false).(bool)
		autoDelete := getConfigBool("auto-delete", queueConfig).OrElse(false).(bool)
		exclusive := getConfigBool("exclusive", queueConfig).OrElse(false).(bool)
		noWait := getConfigBool("no-wait", queueConfig).OrElse(false).(bool)

		queues[queueKey] = &QueueConfiguration{queue, durable, autoDelete, exclusive, noWait}
	}

	return &Configuration{host, port, serverQueue, queues}, nil
}
