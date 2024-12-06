package pkg

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaMessage struct {
	TopicPartition kafka.TopicPartition
	Value          []byte
	Key            []byte
	Timestamp      time.Time
	TimestampType  kafka.TimestampType
	Opaque         interface{}
	Headers        []kafka.Header
}

type KafkaHeader struct {
	Key   string
	Value []byte
}

type ProducerConfig struct {
	TopicName                  string // Topic name to publish messages to
	Acks                       string // Acknowledgment level for the message (e.g., "all", "1", "0")
	Retries                    int    // Number of retries for message delivery
	BatchSize                  int    // Batch size for accumulating messages before sending
	LingerMs                   int    // Wait time for collecting messages before sending (in milliseconds)
	CompressionType            string // Compression type for the messages (e.g., "gzip", "snappy")
	MaxInFlightRequestsPerConn int    // Maximum number of requests allowed per connection at the same time
	ClientID                   string // Kafka producer client ID
	DeliveryTimeoutMs          int    // Delivery timeout in milliseconds
}

type ConsumerConfig struct {
	TopicName              string                    // Topic name to subscribe to
	HandlerFunc            func(*KafkaMessage) error // Function to handle incoming messages
	GroupID                string                    // Consumer group ID
	SessionTimeoutMs       int                       // Session timeout in milliseconds
	HeartbeatIntervalMs    int                       // Interval to send heartbeat messages to Kafka
	EnableAutoOffsetStore  bool                      // Whether to prevent auto-committing offsets
	IsolationLevel         string                    // Message isolation level
	FetchMinBytes          int                       // Minimum number of bytes to fetch per poll
	MaxPartitionFetchBytes int                       // Maximum number of bytes to fetch from each partition
	ClientID               string                    // Kafka consumer client ID
}
