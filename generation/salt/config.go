package salt

import (
	"io/ioutil"
	"github.com/liyinhgqw/typesafe-config/parse"
)

// QueueConfiguration describes how a specific queue is to be created
type QueueConfiguration struct {
	queue     string
	durable   bool
	deletable bool
	exclusive bool
	nowait    bool
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

// ReadConfiguration reads a typesafe-config configuration file that sets up our salt connection
func ReadConfiguration (filename string) (*Configuration, error) {
	configString, err := ioutil.ReadFile(filename)
	if (err != nil) {
		panic(err)
	}

	configMap, err := parse.Parse("salt", string(configString))
	if (err != nil) {
		panic(err)
	}

	config := configMap.GetConfig()
	
	host := getConfigString("rabbitmq.host", config).OrElse("localhost").(string)
	port := getConfigInt("rabbitmq.port", config).OrElse(5672).(int64)
	serverQueue := getConfigString("communications.queue", config).OrElse("salt").(string)
	queues := make(map[string]*QueueConfiguration)
	queueConfigs := getConfigArray("rabbitmq.queues", config).OrElse([0]*parse.Config{}).([]*parse.Config)
	for _, queueConfig := range queueConfigs {
		name := queueConfig.String()
		queueInterface, err := getConfigString("name", queueConfig).Value()
		if nil != err {
			panic(err)
		}
		queue := queueInterface.(string)
		durable := getConfigBool("durable", queueConfig).OrElse(false).(bool)
		autoDelete := getConfigBool("auto-delete", queueConfig).OrElse(false).(bool)
		exclusive := getConfigBool("exclusive", queueConfig).OrElse(false).(bool)
		noWait := getConfigBool("no-wait", queueConfig).OrElse(false).(bool)

		queues[name] = &QueueConfiguration{queue, durable, autoDelete, exclusive, noWait}
	}

	return &Configuration{host, port, serverQueue, queues}, nil
}
