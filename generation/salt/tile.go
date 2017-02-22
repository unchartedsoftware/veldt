package salt

// Tile represents a salt tile type
type Tile struct {
	host string // The hostname of the machine on which the RabbitMQ server resides
	port string // The port to which the RabbitMQ server is listening
}
