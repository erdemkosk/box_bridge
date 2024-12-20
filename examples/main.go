package main

import (
	"errors"
	"log"
	"time"

	"github.com/erdemkosk/box_bridge"
	"github.com/erdemkosk/box_bridge/pkg"
	"github.com/google/uuid"
)

type Foo struct {
	Name string `json:"name"`
	Note string `json:"note"`
}

func main() {

	boxBridge := box_bridge.NewBoxBridge(pkg.NewConfigBuilder().
		WithMongoDBURL("mongodb://localhost:27017").
		WithKafkaURL("localhost:9092").
		WithOutboxCollection("outbox").
		WithInboxCollection("inbox").
		WithRetryAttempts(3).
		Build())

	producerConfig := pkg.ProducerConfig{
		TopicName: "my-topic",
		ClientID:  "my-producer",
	}

	boxBridge.AddProducer(producerConfig)

	boxBridge.Produce(producerConfig, "key", Foo{
		Name: "Erdem Köşk",
		Note: "Hey , box-bridge is amazing mate!",
	}, uuid.New().String()+"-example service", nil)

	// Create handler function for consumer each or one if u return nil it will commited if u return error it wont commited
	handlerFunc := func(msg *pkg.KafkaMessage) error {
		log.Printf("Received message: %s", string(msg.Value))

		err := errors.New("something went wrong , it should not commit any offset!")

		return err
	}

	consumerConfig := pkg.ConsumerConfig{
		TopicName:   "my-topic",
		GroupID:     "my-consumer-group",
		HandlerFunc: handlerFunc,
	}

	boxBridge.AddConsumer(consumerConfig)

	select {
	case <-time.After(10 * time.Second):
		log.Println("Shutting down Kafka manager")
		boxBridge.Shutdown()
	}

}
