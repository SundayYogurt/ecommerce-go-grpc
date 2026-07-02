package queue

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/sundayyogurt/order_service/internal/interfaces"
)

type KafkaConsumer struct {
	Reader      *kafka.Reader
	Handler     interfaces.ConsumerHandler
	ServiceName string
}

func NewKafkaConsumer(broker, topic string, groupID string, handler interfaces.ConsumerHandler) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaConsumer{
		Reader:      reader,
		Handler:     handler,
		ServiceName: "User Service",
	}
}

func (kc *KafkaConsumer) Listen() {

	//listen for messages continuously
	for {
		msg, err := kc.Reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error on reading message: %v", err)
			continue
		}

		log.Printf("Received message: %s", string(msg.Value))

		if err := kc.Handler.HandleMessage(string(msg.Value)); err != nil {
			log.Printf("Error on processing message on handler: %v", err)
		}
	}
}
