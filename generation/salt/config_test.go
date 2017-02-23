package salt



import (
	"testing"
	"reflect"
)



func assertEquals (t *testing.T, name string, expected, actual interface{}) bool {
	if expected != actual {
		var eType = reflect.TypeOf(expected)
		var aType = reflect.TypeOf(actual)
		t.Errorf("%s: Expected %v[%v], got %v[%v]", name, expected, eType, actual, aType)
		return false
	}
	return true
}
func TestFullSaltConfiguration (t *testing.T) {
	config, err := ReadConfiguration("testdata/salt-configuration-full.conf")
	if err != nil {
		t.Error(err)
	}

	assertEquals(t, "Host", "rabbitmq.uncharted.software", config.host)
	assertEquals(t, "Port", int64(1234), config.port)
	assertEquals(t, "Queue", "salt-test-queue", config.serverQueue)
	assertEquals(t, "Queues", 3, len(config.queueConfigurations))
	assertEquals(t, "Queue bunny name", "bunny-queue", config.queueConfigurations["bunny"].queue)
	assertEquals(t, "Queue bunny durability",   false, config.queueConfigurations["bunny"].durable)
	assertEquals(t, "Queue bunny deletability", true,  config.queueConfigurations["bunny"].deletable)
	assertEquals(t, "Queue bunny exclusivity",  false, config.queueConfigurations["bunny"].exclusive)
	assertEquals(t, "Queue bunny no-wait",      true,  config.queueConfigurations["bunny"].noWait)
	assertEquals(t, "Queue lapin name", "lapin-queue", config.queueConfigurations["lapin"].queue)
	assertEquals(t, "Queue lapin durability",   true,  config.queueConfigurations["lapin"].durable)
	assertEquals(t, "Queue lapin deletability", true,  config.queueConfigurations["lapin"].deletable)
	assertEquals(t, "Queue lapin exclusivity",  false, config.queueConfigurations["lapin"].exclusive)
	assertEquals(t, "Queue lapin no-wait",      false, config.queueConfigurations["lapin"].noWait)
	assertEquals(t, "Queue hare name", "hare-queue",   config.queueConfigurations["hare"].queue)
	assertEquals(t, "Queue hare durability",   false,  config.queueConfigurations["hare"].durable)
	assertEquals(t, "Queue hare deletability", false,  config.queueConfigurations["hare"].deletable)
	assertEquals(t, "Queue hare exclusivity",  true,   config.queueConfigurations["hare"].exclusive)
	assertEquals(t, "Queue hare no-wait",      true,   config.queueConfigurations["hare"].noWait)
}

func TestDefaultSaltConfiguration (t *testing.T) {
	config, err := ReadConfiguration("testdata/salt-configuration-empty.conf")
	if err != nil {
		t.Error(err)
	}

	assertEquals(t, "Host", "localhost", config.host)
	assertEquals(t, "Port", int64(5672), config.port)
	assertEquals(t, "Queue", "salt", config.serverQueue)
	assertEquals(t, "Queues", 0, len(config.queueConfigurations))
}
		
