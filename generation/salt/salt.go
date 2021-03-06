package salt

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/streadway/amqp"
)

// This file contains the basic facilities for connecting to and communicating
// with a Salt-based tile server through RabbitMQ
//
// Besides basic setup and configuration, users are meant to communicate with
// the server through three functions:
//
//   * Dataset
//   * QueryTiles
//   * QueryMetadata
//
// Unexported functions herin (especially sendServerMessage) should be
// considered private to this class, not just to the package.

// RabbitMQConnection describes a connection to a RabbitMQ server
type RabbitMQConnection struct {
	connection  *amqp.Connection
	channel     *amqp.Channel
	queues      map[string]amqp.Queue
	serverQueue string
}

var (
	mutex            = sync.Mutex{}
	connections      = make(map[string]*RabbitMQConnection)
	responseChannels = make(map[string]chan<- amqp.Delivery)
	nextMessage      = 0
	emptyResponse    = make([]byte, 0)
)

// NewConnection returns a connection to the Salt tile server via RabbitMQ
func NewConnection(config *Config) (*RabbitMQConnection, error) {
	mutex.Lock()
	rmq, contained := connections[config.Key()]
	if !contained {
		Infof("New connection request")

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

		rmq = &RabbitMQConnection{connection, channel, make(map[string]amqp.Queue), config.serverQueue}

		// Register our standard queues
		for k, v := range config.queueConfigs {
			rmq.Declare(k, v)
		}

		// Start up a consumer on our response channel
		responseQ, err := rmq.GetQueue("response")
		if err != nil {
			return nil, err
		}
		responses, err := channel.Consume(responseQ.Name, "", true, false, false, false, nil)
		if err != nil {
			return nil, err
		}
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

	Infof("connection request fulfilled: %v", rmq)
	return rmq, nil
}

// Close closes this RabbitMQ connection
func (rmq *RabbitMQConnection) Close() {
	Infof("Closing connection")
	rmq.channel.Close()
	rmq.connection.Close()
}

// Declare declares a queue using this RabbitMQ connection
func (rmq *RabbitMQConnection) Declare(qName string, qc *QueueConfig) error {
	q, err := rmq.channel.QueueDeclare(qc.queue, qc.durable, qc.deletable, qc.exclusive, qc.noWait, nil)
	if err != nil {
		return err
	}
	rmq.queues[qName] = q
	Infof("Declaring channel '%s'=(%v)", qName, q)
	return nil
}

// GetQueue gets a predefined channel associated with this connnection by
// canonical name.  If there is no current channel with the given canonical
// name, a temporary channel is created.
func (rmq *RabbitMQConnection) GetQueue(cannonicalName string) (amqp.Queue, error) {
	var err error
	q, contained := rmq.queues[cannonicalName]
	if !contained {
		q, err = rmq.channel.QueueDeclare("", false, true, true, false, nil)
		if err == nil {
			rmq.queues[cannonicalName] = q
		}
	}
	return q, err
}

func nextMessageID() string {
	mutex.Lock()
	defer mutex.Unlock()

	msgID := fmt.Sprintf("message:%d", nextMessage)
	nextMessage++

	return msgID
}

// Dataset sets up a dataset on the Salt server for future use
func (rmq *RabbitMQConnection) Dataset(message []byte) ([]byte, error) {
	return rmq.sendServerMessage("dataset", message)
}

// QueryTiles queries the salt server for a tile
func (rmq *RabbitMQConnection) QueryTiles(message []byte) ([]byte, error) {
	return rmq.sendServerMessage("tiles", message)
}

// QueryMetadata queries the salt server for metadata on a dataset
func (rmq *RabbitMQConnection) QueryMetadata(message []byte) ([]byte, error) {
	return rmq.sendServerMessage("metadata", message)
}

// sendServerMessage is a low-level generic function to do exactly what it says.  It is used by
// Query and Dataset
func (rmq *RabbitMQConnection) sendServerMessage(messageType string, message []byte) ([]byte, error) {
	queryQ, err := rmq.GetQueue(rmq.serverQueue)
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

	Debugf("Publishing message \"%s\"\n\t(query queue: %s(=%s))\n\t(response queue: %s(=%s))\n\t(type: %s)",
		string(message), rmq.serverQueue, queryQ.Name, "response", responseQ.Name, messageType)

	rmq.channel.Publish("", queryQ.Name, false, false,
		amqp.Publishing{
			Type:      messageType,
			Body:      message,
			ReplyTo:   responseQ.Name,
			MessageId: msgID})

	response := <-responseChannel
	Debugf("Response received: \"%s\"", string(response.Body))
	if "error" == response.Type {
		return nil, fmt.Errorf(string(response.Body))
	}

	return response.Body, nil
}
