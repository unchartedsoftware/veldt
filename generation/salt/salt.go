package salt

import (
	"runtime"
	"fmt"
	"sync"
	"github.com/streadway/amqp"
)

// RabbitMQConnection describes a connection to a RabbitMQ server
type RabbitMQConnection struct {
	connection	*amqp.Connection
	channel		*amqp.Channel
	queues		map[string]amqp.Queue
}

var (
	mutex = sync.Mutex{}
	connections = make(map[string]*RabbitMQConnection)
	responseChannels = make(map[string]chan<- amqp.Delivery)
	nextMessage = 0
	emptyResponse = make([]byte, 0)
)

// NewConnection returns a connection to the Salt tile server via RabbitMQ
func NewConnection (config Configuration) (*RabbitMQConnection, error) {
	mutex.Lock()
	rmq, contained := connections[config.Key()]
	if !contained {
		url := fmt.Sprintf("amqp://%s:%d", config.host, config.port)
		var connection *amqp.Connection
		var channel *amqp.Channel

		connection, err := amqp.Dial(url)
		if err != nil {
			return nil, err
		}

		channel, err = connection.Channel()
		if err != nil {
			return nil, err
		}

		rmq = &RabbitMQConnection{connection, channel, make(map[string]amqp.Queue)}

		// Register our standard queues
		for k, v := range config.queueConfigurations {
			rmq.Declare(k, v)
		}

		// Start up a consumer on our response channel
		responseQ, err := rmq.GetQueue("response")
		if err != nil {
			return nil, err
		}
		responses, err := channel.Consume(responseQ.Name, "", true, false, false, false, nil)
		go func() {
			for response := range responses {
				msgID := response.MessageId
				responseChannel := responseChannels[msgID]
				delete(responseChannels, msgID)
				responseChannel <- response
			}
		}()

		// Store this connection for later reuse
		connections[config.Key()] = rmq
	}
	mutex.Unlock()
	runtime.Gosched()

	return rmq, nil
}

// Close closes this RabbitMQ connection
func (rmq *RabbitMQConnection) Close () {
	rmq.channel.Close()
	rmq.connection.Close()
}




// Declare declares a queue using this RabbitMQ connection
func (rmq *RabbitMQConnection) Declare (qName string, qc *QueueConfiguration) error {
	q, err := rmq.channel.QueueDeclare(qc.queue, qc.durable, qc.deletable, qc.exclusive, qc.noWait, nil)
	if (err != nil) {
		return err
	}
	rmq.queues[qName] = q
	return nil
}

// GetQueue gets a predefined channel associated with this connnection by
// cannonical name.  If there is no current channel with the given cannonical
// name, a temporary channel is created.
func (rmq *RabbitMQConnection) GetQueue (cannonicalName string) (amqp.Queue, error) {
	var err error
	q, contained := rmq.queues[cannonicalName]
	if (!contained) {
		q, err = rmq.channel.QueueDeclare("", false, true, true, false, nil)
		if err == nil {
			rmq.queues[cannonicalName] = q
		}
	}
	return q, err
}

func nextMessageID () string {
	mutex.Lock()
	defer mutex.Unlock()

	msgID := fmt.Sprintf("message:%d", nextMessage)
	nextMessage++

	return msgID
}

// Query the salt server for a tile
func (rmq *RabbitMQConnection) Query (queue string, message []byte) ([]byte, error) {
	queryQ, err := rmq.GetQueue(queue)
	if err != nil {
		return emptyResponse, err
	}
	responseQ, err := rmq.GetQueue("response")
	if err != nil {
		return emptyResponse, err
	}

	msgID := nextMessageID()
	responseChannel := make(chan amqp.Delivery)
	responseChannels[msgID] = responseChannel

	rmq.channel.Publish("", queryQ.Name, false, false, amqp.Publishing{Body: message, ReplyTo: responseQ.Name, MessageId: msgID})

	response := <- responseChannel
	return response.Body, nil
}
