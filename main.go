package main

import (
	"log"
	"time"

	"github.com/erdemkosk/box_bridge/pkg"
	"github.com/erdemkosk/box_bridge/pkg/model"
)

type Foo struct {
	Name string `json:"name"`
	Note string `json:"note"`
}

func main() {
	boxBridge := pkg.NewBoxBridge(pkg.NewConfig("localhost:27017", "", "", "", 0))

	producerConfig := model.ProducerConfig{
		TopicName: "my-topic",
		ClientID:  "my-producer",
	}

	boxBridge.AddProducer(producerConfig)

	boxBridge.Produce(producerConfig, "key1", Foo{
		Name: "Erdem Köşk",
		Note: "Hey , box-bridge is amazing mate!",
	})

	// Create handler function for consumer each or one
	handlerFunc := func(msg *model.KafkaMessage) error {
		log.Printf("Received message: %s", string(msg.Value))
		return nil
	}

	consumerConfig := model.ConsumerConfig{
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
