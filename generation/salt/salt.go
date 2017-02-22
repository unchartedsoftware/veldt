package salt

import (
	"github.com/streadway/amqp"
)

// RabbitMQConnection describes a connection to a RabbitMQ server
type RabbitMQConnection struct {
	connection	*amqp.Connection
	channel		*amqp.Channel
	queues		map[string]amqp.Queue
}

type queueDefinition struct {
	durable bool
	deletable bool
	exclusive bool
	nowait bool
}
type queueRedefinition func (*queueDefinition)

const (
)

var (
)

// NewConnection returns a connection to the Salt tile server via RabbitMQ
func NewConnection (host, port string) (RabbitMQConnection, error) {
	var err error
	url := "amqp://"+host+":"+port
	var connection *amqp.Connection
	var channel *amqp.Channel
	connection, err = amqp.Dial(url)
	if err != nil {
		channel, err = connection.Channel()
	}
	return RabbitMQConnection{connection, channel, make(map[string]amqp.Queue)}, err
}

// Close closes this RabbitMQ connection
func (rmq *RabbitMQConnection) Close () {
	rmq.channel.Close()
	rmq.connection.Close()
}


// Functions that can be passed into queue definition to alter the way the queue is made

// Durable is used to make a durable queue
func Durable (qd *queueDefinition) {
	qd.durable = true
}

// Deletable is used to create a deletable queue
func Deletable (qd *queueDefinition) {
	qd.deletable = true
}

// Exclusive is used to create an exclusive queue
func Exclusive (qd *queueDefinition) {
	qd.exclusive = true
}

// NoWait is used to create a no-wait queue
func NoWait (qd *queueDefinition) {
	qd.nowait = true;
}


// Declare declares a queue using this RabbitMQ connection
func (rmq *RabbitMQConnection) Declare (queue string, options ...queueRedefinition) error {
	qd := &queueDefinition{false, false, false, false}
	for _, opt := range options {
		opt(qd)
	}
	q, err := rmq.channel.QueueDeclare(queue, qd.durable, qd.deletable, qd.exclusive, qd.nowait, nil)
	if (err != nil) {
		return err
	}
	rmq.queues[queue] = q
	return nil
}
