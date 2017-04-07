package salt

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Salt configuration", func() {
	It("should read all properties correctly from a test config file", func() {
		config, err := ReadConfig("testdata/salt-configuration-full.conf")
		Expect(err).To(BeNil())

		Expect(config.host).To(Equal("rabbitmq.uncharted.software"), "Host parameter")
		Expect(config.port).To(Equal(int64(1234)), "Port")
		Expect(config.serverQueue).To(Equal("salt-test-queue"), "Queue")
		Expect(len(config.queueConfigs)).To(Equal(3), "Queues")
		Expect(config.queueConfigs["bunny"].queue).To(Equal("bunny-queue"), "Queue bunny name")
		Expect(config.queueConfigs["bunny"].durable).To(Equal(false), "Queue bunny durability")
		Expect(config.queueConfigs["bunny"].deletable).To(Equal(true), "Queue bunny deletability")
		Expect(config.queueConfigs["bunny"].exclusive).To(Equal(false), "Queue bunny exclusivity")
		Expect(config.queueConfigs["bunny"].noWait).To(Equal(true), "Queue bunny no-wait")
		Expect(config.queueConfigs["lapin"].queue).To(Equal("lapin-queue"), "Queue lapin name")
		Expect(config.queueConfigs["lapin"].durable).To(Equal(true), "Queue lapin durability")
		Expect(config.queueConfigs["lapin"].deletable).To(Equal(true), "Queue lapin deletability")
		Expect(config.queueConfigs["lapin"].exclusive).To(Equal(false), "Queue lapin exclusivity")
		Expect(config.queueConfigs["lapin"].noWait).To(Equal(false), "Queue lapin no-wait")
		Expect(config.queueConfigs["hare"].queue).To(Equal("hare-queue"), "Queue hare name")
		Expect(config.queueConfigs["hare"].durable).To(Equal(false), "Queue hare durability")
		Expect(config.queueConfigs["hare"].deletable).To(Equal(false), "Queue hare deletability")
		Expect(config.queueConfigs["hare"].exclusive).To(Equal(true), "Queue hare exclusivity")
		Expect(config.queueConfigs["hare"].noWait).To(Equal(true), "Queue hare no-wait")
	})

	It("Should default all configuration parameters properly", func() {
		config, err := ReadConfig("testdata/salt-configuration-empty.conf")
		Expect(err).To(BeNil())

		Expect(config.host).To(Equal("localhost"), "Host")
		Expect(config.port).To(Equal(int64(5672)), "Port")
		Expect(config.serverQueue).To(Equal("salt"), "Queue")
		Expect(len(config.queueConfigs)).To(Equal(0), "Queues")
	})
})
