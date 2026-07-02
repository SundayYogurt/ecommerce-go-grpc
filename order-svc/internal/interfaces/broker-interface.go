package interfaces

// recive message from queue
type ConsumerHandler interface {
	HandleMessage(message string) error
}

// send message to queue
type ProducerHandler interface {
	PublishMessage(key, value []byte) error
}
