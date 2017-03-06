package salt

import (
	"runtime"
	"fmt"
	"sync"

	"github.com/streadway/amqp"
	"github.com/unchartedsoftware/plog"
)

// RabbitMQConnection describes a connection to a RabbitMQ server
type RabbitMQConnection struct {
	connection	*amqp.Connection
	channel		*amqp.Channel
	queues		map[string]amqp.Queue
	serverQueue string
}

const (
	// Yellow "SALT" - code derived from github.com/mgutz/ansi,
	// but I wanted it as a const, so couldn't use that directly
	preLog = "\033[1;38;5;3mSALT\033[0m: "
	// And codes to make the message red, similarly as constants
	preMsg = "\033[1;97;3m"
	postMsg = "\033[0m"
)

var (
	mutex = sync.Mutex{}
	connections = make(map[string]*RabbitMQConnection)
	responseChannels = make(map[string]chan<- amqp.Delivery)
	nextMessage = 0
	emptyResponse = make([]byte, 0)
)

// NewConnection returns a connection to the Salt tile server via RabbitMQ
func NewConnection (config *Configuration) (*RabbitMQConnection, error) {
	mutex.Lock()
	rmq, contained := connections[config.Key()]
	if !contained {
		log.Infof(preLog+"New connection request")

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

	log.Infof(preLog+"connection request fulfilled: %v", rmq)
	return rmq, nil
}

// Close closes this RabbitMQ connection
func (rmq *RabbitMQConnection) Close () {
	log.Infof(preLog+"Closing connection")
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
	log.Infof(preLog+"Declaring channel '%s'=(%v)", qName, q)
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

// Dataset sets up a dataset on the Salt server for future use
func (rmq *RabbitMQConnection) Dataset (message []byte) ([]byte, error) {
	return rmq.sendServerMessage("dataset", message)
}

// QueryTiles queries the salt server for a tile
func (rmq *RabbitMQConnection) QueryTiles (message []byte) ([]byte, error) {
	return rmq.sendServerMessage("tiles", message)
}

// QueryMetadata queries the salt server for metadata on a dataset
func (rmq *RabbitMQConnection) QueryMetadata (message []byte) ([]byte, error) {
	return rmq.sendServerMessage("metadata", message)
}

// sendServerMessage is a low-level generic function to do exactly what it says.  It is used by
// Query and Dataset
func (rmq *RabbitMQConnection) sendServerMessage (messageType string, message []byte) ([]byte, error) {
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

	log.Infof(preLog+"Publishing message \"%s%s%s\" (query queue: %s(=%s)) (response queue: %s(=%s))",
		preMsg, string(message), postMsg, rmq.serverQueue, queryQ.Name, "response", responseQ.Name)

	rmq.channel.Publish("", queryQ.Name, false, false,
		amqp.Publishing{
			Type: messageType,
			Body: message,
			ReplyTo: responseQ.Name,
			MessageId: msgID})

	response := <- responseChannel
	log.Infof(preLog+"Response received: \"%s%s%s\"", preMsg, string(response.Body), postMsg)
	return response.Body, nil
}
