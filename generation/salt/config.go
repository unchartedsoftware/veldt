package salt

import (
	"fmt"
	"github.com/liyinhgqw/typesafe-config/parse"
	"io/ioutil"
	"sort"
	"strings"
)

// This file describes the configuration information needed to connect to a
// salt tile server.  The Config class is the main class herein.
//
// Also included are utility functions to read a configuration from a file -
// ReadConfig - and to read a dataset configuration (currently just
// stored as a string, as only the server has a real need to parse it) in
// ReadDatasetConfig.  These two functions are intended for use in the
// main routine of an application, to read config objects from file, so they
// can be passed into salt tile constructors.

// QueueConfig describes how a specific queue is to be created
type QueueConfig struct {
	queue     string
	durable   bool
	deletable bool
	exclusive bool
	noWait    bool
}

// Config describes how the salt connection is to be created
type Config struct {
	host         string
	port         int64
	queueConfigs map[string]*QueueConfig
	serverQueue  string
}

func tOrF(b bool) string {
	if b {
		return "T"
	}
	return "F"
}

// Key returns a unique key that completely describes this configuration
func (c *Config) Key() string {
	qcs := make([]string, len(c.queueConfigs))
	i := 0
	for k := range c.queueConfigs {
		qcs[i] = k
	}
	sort.Strings(qcs)
	qcKey := ""
	first := true
	for _, k := range qcs {
		qc := c.queueConfigs[k]
		if first {
			qcKey = fmt.Sprintf("%s.%s.%s%s%s%s", k, qc.queue, tOrF(qc.durable), tOrF(qc.deletable), tOrF(qc.exclusive), tOrF(qc.noWait))
		} else {
			qcKey = qcKey + fmt.Sprintf(":%s.%s.%s%s%s%s", k, qc.queue, tOrF(qc.durable), tOrF(qc.deletable), tOrF(qc.exclusive), tOrF(qc.noWait))
		}
	}

	return fmt.Sprintf("%s:%d|%s|%s", c.host, c.port, c.serverQueue, qcKey)
}

func getConfigString(key string, config *parse.Config) Option {
	result, err := config.GetString(key)
	if err != nil {
		return Option{false, nil}
	}
	return Option{true, result}
}
func getConfigInt(key string, config *parse.Config) Option {
	result, err := config.GetInt(key)
	if err != nil {
		return Option{false, nil}
	}
	return Option{true, result}
}
func getConfigBool(key string, config *parse.Config) Option {
	result, err := config.GetBool(key)
	if err != nil {
		return Option{false, nil}
	}
	return Option{true, result}
}
func getConfigArray(key string, config *parse.Config) Option {
	result, err := config.GetBool(key)
	if err != nil {
		return Option{false, nil}
	}
	return Option{true, result}
}
func getSubConfig(key string, config *parse.Config) Option {
	result, err := config.GetValue(key)
	if err != nil {
		return Option{false, nil}
	}
	return Option{true, result}
}

func stripTerminalQuotes(text string) string {
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
func getKeys(root parse.Node, path ...string) ([]string, error) {
	if root.Type() != parse.NodeMap {
		return nil, fmt.Errorf("root node is not a map")
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
		return nil, fmt.Errorf("node %s not found in map", path[0])
	}
	return getKeys(elem, path[1:]...)
}

// ReadDatasetConfig reads a file into a string, so it can be passed to Salt
func ReadDatasetConfig(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// ReadConfig reads a typesafe-config configuration file that sets up our salt connection
func ReadConfig(filename string) (*Config, error) {
	// Read the config file
	configString, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Parse it out into a config map
	configMap, err := parse.Parse("salt", string(configString))
	if err != nil {
		return nil, err
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
	queues := make(map[string]*QueueConfig)
	for _, queueKey := range queueKeys {
		queueConfig := getSubConfig("rabbitmq.queues."+queueKey, config).OrElse(&parse.Config{}).(*parse.Config)
		queueInterface, err := getConfigString("name", queueConfig).Value()
		if err != nil {
			return nil, err
		}
		queue := stripTerminalQuotes(queueInterface.(string))
		durable := getConfigBool("durable", queueConfig).OrElse(false).(bool)
		autoDelete := getConfigBool("auto-delete", queueConfig).OrElse(false).(bool)
		exclusive := getConfigBool("exclusive", queueConfig).OrElse(false).(bool)
		noWait := getConfigBool("no-wait", queueConfig).OrElse(false).(bool)

		queues[queueKey] = &QueueConfig{queue, durable, autoDelete, exclusive, noWait}
	}

	return &Config{host, port, queues, serverQueue}, nil
}
