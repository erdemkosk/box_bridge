package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/erdemkosk/box_bridge/pkg"
)

type KafkaManager struct {
	producers map[string]*kafka.Producer
	consumers map[string]*kafka.Consumer
	brokers   string
	mu        sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewKafkaManager(brokers string) *KafkaManager {
	ctx, cancel := context.WithCancel(context.Background())

	return &KafkaManager{
		brokers:   brokers,
		producers: make(map[string]*kafka.Producer),
		consumers: make(map[string]*kafka.Consumer),
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (k *KafkaManager) InitProducer(config pkg.ProducerConfig) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	producerConfig := &kafka.ConfigMap{
		"bootstrap.servers":                     k.brokers,
		"acks":                                  defaultString(config.Acks, "all"),
		"retries":                               defaultInt(config.Retries, 5),
		"batch.size":                            defaultInt(config.BatchSize, 16384),
		"linger.ms":                             defaultInt(config.LingerMs, 5),
		"compression.type":                      defaultString(config.CompressionType, "none"),
		"max.in.flight.requests.per.connection": defaultInt(config.MaxInFlightRequestsPerConn, 5),
		"client.id":                             defaultString(config.ClientID, "default-producer"),
		"delivery.timeout.ms":                   defaultInt(config.DeliveryTimeoutMs, 30000),
	}

	producer, err := kafka.NewProducer(producerConfig)
	if err != nil {
		return err
	}

	k.producers[config.TopicName] = producer
	log.Printf("BOX-BRIDGE: Producer for topic %s initialized", config.TopicName)
	return nil
}

func (k *KafkaManager) StartConsumer(topicConfig pkg.ConsumerConfig) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	consumerConfig := &kafka.ConfigMap{
		"bootstrap.servers":         k.brokers,
		"group.id":                  topicConfig.GroupID,
		"auto.offset.reset":         "earliest",
		"enable.auto.commit":        false,
		"session.timeout.ms":        defaultInt(topicConfig.SessionTimeoutMs, 10000),
		"heartbeat.interval.ms":     defaultInt(topicConfig.HeartbeatIntervalMs, 3000),
		"enable.auto.offset.store":  topicConfig.EnableAutoOffsetStore,
		"isolation.level":           defaultString(topicConfig.IsolationLevel, "read_committed"),
		"fetch.min.bytes":           defaultInt(topicConfig.FetchMinBytes, 1),
		"max.partition.fetch.bytes": defaultInt(topicConfig.MaxPartitionFetchBytes, 1048576),
		"client.id":                 defaultString(topicConfig.ClientID, "default-client"),
	}

	consumer, err := kafka.NewConsumer(consumerConfig)
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

				convertedValue := pkg.KafkaMessage{
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

func (k *KafkaManager) Produce(topic string, key []byte, value []byte, headers []pkg.KafkaHeader) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	var convertedHeaders []kafka.Header
	if headers != nil && len(headers) > 0 {
		convertedHeaders = ConvertKafkaHeadersToWrapperHeaders(headers)
	}

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
		Key:     key,
		Value:   value,
		Headers: convertedHeaders,
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

func (k *KafkaManager) CommitOffset(msg *pkg.KafkaMessage) error {
	offset := msg.TopicPartition.Offset
	log.Printf("Attempting to commit offset %v for message %v", offset, msg.TopicPartition)

	_, err := k.consumers[*msg.TopicPartition.Topic].CommitOffsets([]kafka.TopicPartition{
		{
			Topic:     msg.TopicPartition.Topic,
			Partition: msg.TopicPartition.Partition,
			Offset:    offset + 1,
		},
	})
	if err != nil {
		log.Printf("Error committing offset %v for message %v: %v", offset, msg.TopicPartition, err)
		return err
	}

	log.Printf("Successfully committed offset %v for message %v", offset, msg.TopicPartition)
	return nil
}
