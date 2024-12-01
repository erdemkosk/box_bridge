package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/erdemkosk/box_bridge/internal/db"
	"github.com/erdemkosk/box_bridge/internal/db/models"
	"github.com/erdemkosk/box_bridge/internal/kafka"
	kafkaModels "github.com/erdemkosk/box_bridge/pkg/model"
)

type Boxbridge struct {
	config *BoxBridgeConfig
}

var boxbridgeRunner *Boxbridge
var kafkaManager *kafka.KafkaManager
var mongoManager *db.MongoManager

func NewBoxBridge(config *BoxBridgeConfig) *Boxbridge {
	var mongoErr error

	mongoManager, mongoErr = db.NewMongoManager(config.MongoDBURL)

	if mongoErr != nil {
		panic(mongoErr)
	}

	kafkaManager = kafka.NewKafkaManager(config.KafkaURL, "my-consumer-group")

	return boxbridgeRunner
}

func (bb *Boxbridge) AddProducer(producerConfig kafkaModels.ProducerConfig) {
	if err := kafkaManager.InitProducer(producerConfig); err != nil {
		log.Fatalf("Error initializing producer: %v", err)
	}
}

// (Kafka → Wrapper -> Inbox -> Hander)
func (bb *Boxbridge) AddConsumer(consumerConfig kafkaModels.ConsumerConfig) {

	wrappedHandler := func(msg *kafkaModels.KafkaMessage) error {

		inboxMessage := models.Inbox{
			Offset:    fmt.Sprintf("%v", msg.TopicPartition.Offset),
			Key:       string(msg.Key),
			Content:   string(msg.Value),
			Status:    "Received",
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}

		err := mongoManager.SaveToInbox(inboxMessage)
		if err != nil {
			log.Printf("Error saving to inbox: %v", err)
		}

		return consumerConfig.HandlerFunc(msg)
	}

	if err := kafkaManager.StartConsumer(kafkaModels.ConsumerConfig{
		TopicName:   consumerConfig.TopicName,
		GroupID:     consumerConfig.GroupID,
		HandlerFunc: wrappedHandler,
	}); err != nil {
		log.Fatalf("Error starting consumer: %v", err)
	}
}

// (Outbox → Kafka -> Update Status Of Outbox)
func (bb *Boxbridge) Produce(producerConfig kafkaModels.ProducerConfig, key string, message interface{}) error {

	jsonKey, _ := json.Marshal(key)
	jsonMessage, _ := json.Marshal(message)

	err := mongoManager.SaveToOutbox(models.Outbox{
		MessageID:  key,
		Topic:      producerConfig.TopicName,
		Content:    string(jsonMessage),
		Status:     "WaitingForSendingKafka",
		RetryCount: 3,
		CreatedAt:  time.Now().Format(time.RFC3339),
		UpdatedAt:  time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to save to outbox: %v", err)
	}

	if err := kafkaManager.Produce(producerConfig.TopicName, jsonKey, jsonMessage); err != nil {
		log.Printf("Error sending message: %v", err)
	} else {
		log.Println("Message successfully sent to Kafka!")
	}

	err = mongoManager.UpdateOutboxStatus(key, "SentToKafka")
	if err != nil {
		return fmt.Errorf("failed to update outbox status: %v", err)
	}

	return nil
}

func (bb *Boxbridge) Shutdown() {
	kafkaManager.Shutdown()
	mongoManager.Shutdown()
}
