package box_bridge

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/erdemkosk/box_bridge/internal/db"
	"github.com/erdemkosk/box_bridge/internal/kafka"

	"github.com/erdemkosk/box_bridge/pkg"
	"github.com/google/uuid"
)

type Boxbridge struct {
	config *pkg.BoxBridgeConfig
}

var boxbridgeRunner *Boxbridge
var kafkaManager *kafka.KafkaManager
var mongoManager *db.MongoManager

func NewBoxBridge(config *pkg.BoxBridgeConfig) *Boxbridge {
	var mongoErr error

	mongoManager, mongoErr = db.NewMongoManager(config.MongoDBURL)

	if mongoErr != nil {
		panic(mongoErr)
	}

	kafkaManager = kafka.NewKafkaManager(config.KafkaURL)

	return boxbridgeRunner
}

func (bb *Boxbridge) AddProducer(producerConfig pkg.ProducerConfig) {
	if err := kafkaManager.InitProducer(producerConfig); err != nil {
		log.Fatalf("Error initializing producer: %v", err)
	}
}

// (Kafka → Wrapper -> Inbox -> Hander)
func (bb *Boxbridge) AddConsumer(consumerConfig pkg.ConsumerConfig) {

	wrappedHandler := func(msg *pkg.KafkaMessage) error {

		correlationID, found := kafka.GetHeaderValue(msg.Headers, "CorrelationID")
		if !found {
			log.Println("CorrelationID not found in message headers")
			correlationID = "unknown"
		}

		inboxMessage := db.Inbox{
			Topic:         *msg.TopicPartition.Topic,
			CorrelationID: correlationID,
			Offset:        fmt.Sprintf("%v", msg.TopicPartition.Offset),
			Key:           string(msg.Key),
			Content:       string(msg.Value),
			Status:        "Received",
			CreatedAt:     time.Now().Format(time.RFC3339),
			UpdatedAt:     time.Now().Format(time.RFC3339),
		}

		err := mongoManager.SaveToInbox(inboxMessage)
		if err != nil {
			log.Printf("Error saving to inbox: %v", err)
		}

		err = consumerConfig.HandlerFunc(msg)
		if err != nil {
			log.Printf("Error in handler function for message: %v", err)

			err = mongoManager.UpdateInboxStatus(correlationID, "FailedToProcess")
			if err != nil {
				return fmt.Errorf("failed to update inbox status: %v", err)
			}

			return err
		}

		err = kafkaManager.CommitOffset(msg)
		if err != nil {
			log.Printf("Error committing offset for message: %v", err)
			return err
		}

		err = mongoManager.UpdateInboxStatus(correlationID, "Processed")
		if err != nil {
			return fmt.Errorf("failed to update inbox status: %v", err)
		}

		log.Printf("Message successfully processed and offset committed for %v", msg.TopicPartition)

		return nil

	}

	if err := kafkaManager.StartConsumer(pkg.ConsumerConfig{
		TopicName:   consumerConfig.TopicName,
		GroupID:     consumerConfig.GroupID,
		HandlerFunc: wrappedHandler,
	}); err != nil {
		log.Fatalf("Error starting consumer: %v", err)
	}
}

// (Outbox → Kafka -> Update Status Of Outbox)
func (bb *Boxbridge) Produce(producerConfig pkg.ProducerConfig, key string, message interface{}, correlationID string, headers []pkg.KafkaHeader) error {

	if correlationID == "" {
		correlationID = uuid.New().String()
	}

	jsonKey, _ := json.Marshal(key)
	jsonMessage, _ := json.Marshal(message)

	err := mongoManager.SaveToOutbox(db.Outbox{
		Key:           key,
		CorrelationID: correlationID,
		Topic:         producerConfig.TopicName,
		Content:       string(jsonMessage),
		Status:        "WaitingForSendingKafka",
		RetryCount:    3,
		CreatedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to save to outbox: %v", err)
	}

	headers = append(headers, pkg.KafkaHeader{
		Key:   "CorrelationID",
		Value: []byte(correlationID),
	})

	if err := kafkaManager.Produce(producerConfig.TopicName, jsonKey, jsonMessage, headers); err != nil {
		log.Printf("Error sending message: %v", err)
	} else {
		log.Println("Message successfully sent to Kafka!")
	}

	err = mongoManager.UpdateOutboxStatus(correlationID, "SentToKafka")
	if err != nil {
		return fmt.Errorf("failed to update outbox status: %v", err)
	}

	return nil
}

func (bb *Boxbridge) Shutdown() {
	kafkaManager.Shutdown()
	mongoManager.Shutdown()
}
