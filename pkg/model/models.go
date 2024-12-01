package model

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
	AutoCommit             bool                      // Whether offsets should be automatically committed
	MaxPollRecords         int                       // Maximum number of records to fetch in a single poll
	SessionTimeoutMs       int                       // Session timeout in milliseconds (if no heartbeat is received, the session expires)
	HeartbeatIntervalMs    int                       // Interval to send heartbeat messages to Kafka (in milliseconds)
	EnableAutoOffsetStore  bool                      // Whether to prevent auto-committing offsets
	IsolationLevel         string                    // Message isolation level (e.g., "read_committed")
	RebalanceTimeoutMs     int                       // Timeout for rebalance operations (in milliseconds)
	FetchMinBytes          int                       // Minimum number of bytes to fetch per poll
	FetchMaxWaitMs         int                       // Maximum wait time to fetch data (in milliseconds)
	MaxPartitionFetchBytes int                       // Maximum number of bytes to fetch from each partition
	RetryBackoffMs         int                       // Backoff time before retrying after a failure (in milliseconds)
	ClientID               string                    // Kafka consumer client ID
}
