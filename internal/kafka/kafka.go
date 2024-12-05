package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/erdemkosk/box_bridge/pkg/model"
)

const (
	KafkaServer  = "localhost:9092"
	KafkaTopic   = "erdem"
	KafkaGroupId = "product-service"
)

func SendMessageToKafka(topic string, message interface{}) error {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": KafkaServer,
	})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	value, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, nil)

	producer.Flush(1000)

	return err
}

func GetMessageFromKafka(targetType interface{}) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": KafkaServer,
		"group.id":          KafkaGroupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	topic := KafkaTopic
	consumer.SubscribeTopics([]string{topic}, nil)

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			var value interface{}
			err := json.Unmarshal(msg.Value, &value)
			if err != nil {
				fmt.Printf("Error decoding message: %v\n", err)
				continue
			}

			fmt.Printf("Received Order: %+v\n", value)
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func GetHeaderValue(headers []kafka.Header, key string) (string, bool) {
	for _, header := range headers {
		if header.Key == key {
			return string(header.Value), true
		}
	}
	return "", false
}

func ConvertKafkaHeadersToWrapperHeaders(kafkaHeaders []model.KafkaHeader) []kafka.Header {
	headers := make([]kafka.Header, len(kafkaHeaders))
	for i, kh := range kafkaHeaders {
		headers[i] = kafka.Header{
			Key:   kh.Key,
			Value: kh.Value,
		}
	}
	return headers
}
