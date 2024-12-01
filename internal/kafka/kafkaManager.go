package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/erdemkosk/box_bridge/pkg/model"
)

type KafkaManager struct {
	producers map[string]*kafka.Producer
	consumers map[string]*kafka.Consumer
	config    *kafka.ConfigMap
	mu        sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewKafkaManager(brokers string, groupID string) *KafkaManager {
	ctx, cancel := context.WithCancel(context.Background())

	config := &kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}

	return &KafkaManager{
		config:    config,
		producers: make(map[string]*kafka.Producer),
		consumers: make(map[string]*kafka.Consumer),
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (k *KafkaManager) InitProducer(config model.ProducerConfig) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	producer, err := kafka.NewProducer(k.config)
	if err != nil {
		return err
	}

	k.producers[config.TopicName] = producer
	log.Printf("BOX-BRIDGE: Producer for topic %s initialized", config.TopicName)
	return nil
}

func (k *KafkaManager) StartConsumer(topicConfig model.ConsumerConfig) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	consumer, err := kafka.NewConsumer(k.config)
	if err != nil {
		return err
	}

	err = consumer.SubscribeTopics([]string{topicConfig.TopicName}, nil)
	if err != nil {
		return err
	}

	k.consumers[topicConfig.TopicName] = consumer

	go func() {
		log.Printf("BOX-BRIDGE: Consumer for topic %s started", topicConfig.TopicName)
		for {
			select {
			case <-k.ctx.Done():
				log.Printf("BOX-BRIDGE: Consumer for topic %s shutting down", topicConfig.TopicName)
				return
			default:
				msg, err := consumer.ReadMessage(-1)
				convertedValue := model.KafkaMessage{
					TopicPartition: msg.TopicPartition,
					Value:          msg.Value,
					Key:            msg.Key,
					Timestamp:      msg.Timestamp,
					TimestampType:  msg.TimestampType,
					Opaque:         msg.Opaque,
					Headers:        msg.Headers,
				}

				if err != nil {
					log.Printf("BOX-BRIDGE: Error reading message from topic %s: %v", topicConfig.TopicName, err)
					continue
				}
				if err := topicConfig.HandlerFunc(&convertedValue); err != nil {
					log.Printf("BOX-BRIDGE: Error handling message from topic %s: %v", topicConfig.TopicName, err)
				}
			}
		}
	}()
	return nil
}

func (k *KafkaManager) Produce(topic string, key []byte, value []byte) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	producer, ok := k.producers[topic]
	if !ok {
		return logErrorf("BOX-BRIDGE: Producer for topic %s not initialized", topic)
	}

	deliveryChan := make(chan kafka.Event, 1)
	defer close(deliveryChan)

	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   key,
		Value: value,
	}, deliveryChan)
	if err != nil {
		return err
	}

	e := <-deliveryChan
	msg := e.(*kafka.Message)
	if msg.TopicPartition.Error != nil {
		return logErrorf("BOX-BRIDGE: Failed to deliver message: %v", msg.TopicPartition.Error)
	}

	log.Printf("BOX-BRIDGE: Message delivered to %v", msg.TopicPartition)
	return nil
}

func (k *KafkaManager) Shutdown() {
	k.cancel()

	k.mu.Lock()
	defer k.mu.Unlock()

	for topic, consumer := range k.consumers {
		log.Printf("BOX-BRIDGE: Closing consumer for topic %s", topic)
		consumer.Close()
	}

	for topic, producer := range k.producers {
		log.Printf("BOX-BRIDGE: Closing producer for topic %s", topic)
		producer.Close()
	}

	log.Println("BOX-BRIDGE: Kafka Manager shut down gracefully")
}

func logErrorf(format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	log.Println(err)
	return err
}
