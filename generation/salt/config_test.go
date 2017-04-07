package salt

import (
	"reflect"
	"testing"
)

func assertEquals(t *testing.T, name string, expected, actual interface{}) bool {
	if expected != actual {
		var eType = reflect.TypeOf(expected)
		var aType = reflect.TypeOf(actual)
		t.Errorf("%s: Expected %v[%v], got %v[%v]", name, expected, eType, actual, aType)
		return false
	}
	return true
}
func TestFullSaltConfig(t *testing.T) {
	config, err := ReadConfig("testdata/salt-configuration-full.conf")
	if err != nil {
		t.Error(err)
	}

	assertEquals(t, "Host", "rabbitmq.uncharted.software", config.host)
	assertEquals(t, "Port", int64(1234), config.port)
	assertEquals(t, "Queue", "salt-test-queue", config.serverQueue)
	assertEquals(t, "Queues", 3, len(config.queueConfigs))
	assertEquals(t, "Queue bunny name", "bunny-queue", config.queueConfigs["bunny"].queue)
	assertEquals(t, "Queue bunny durability", false, config.queueConfigs["bunny"].durable)
	assertEquals(t, "Queue bunny deletability", true, config.queueConfigs["bunny"].deletable)
	assertEquals(t, "Queue bunny exclusivity", false, config.queueConfigs["bunny"].exclusive)
	assertEquals(t, "Queue bunny no-wait", true, config.queueConfigs["bunny"].noWait)
	assertEquals(t, "Queue lapin name", "lapin-queue", config.queueConfigs["lapin"].queue)
	assertEquals(t, "Queue lapin durability", true, config.queueConfigs["lapin"].durable)
	assertEquals(t, "Queue lapin deletability", true, config.queueConfigs["lapin"].deletable)
	assertEquals(t, "Queue lapin exclusivity", false, config.queueConfigs["lapin"].exclusive)
	assertEquals(t, "Queue lapin no-wait", false, config.queueConfigs["lapin"].noWait)
	assertEquals(t, "Queue hare name", "hare-queue", config.queueConfigs["hare"].queue)
	assertEquals(t, "Queue hare durability", false, config.queueConfigs["hare"].durable)
	assertEquals(t, "Queue hare deletability", false, config.queueConfigs["hare"].deletable)
	assertEquals(t, "Queue hare exclusivity", true, config.queueConfigs["hare"].exclusive)
	assertEquals(t, "Queue hare no-wait", true, config.queueConfigs["hare"].noWait)
}

func TestDefaultSaltConfig(t *testing.T) {
	config, err := ReadConfig("testdata/salt-configuration-empty.conf")
	if err != nil {
		t.Error(err)
	}

	assertEquals(t, "Host", "localhost", config.host)
	assertEquals(t, "Port", int64(5672), config.port)
	assertEquals(t, "Queue", "salt", config.serverQueue)
	assertEquals(t, "Queues", 0, len(config.queueConfigs))
}
