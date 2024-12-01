package pkg

type Config struct {
	MongoDBURL    string `mapstructure:"MONGODB_URL"`
	KafkaURL      string `mapstructure:"KAFKA_URL"`
	OutboxColl    string `mapstructure:"OUTBOX_COLLECTION"`
	InboxColl     string `mapstructure:"INBOX_COLLECTION"`
	RetryAttempts int    `mapstructure:"RETRY_ATTEMPTS"`
}

const (
	defaultMongoDBURL    = "mongodb://localhost:27017"
	defaultKafkaURL      = "localhost:9092"
	defaultOutboxColl    = "outbox"
	defaultInboxColl     = "inbox"
	defaultRetryAttempts = 3
)

func NewConfig(mongoDBURL string, kafkaURL string, outboxColl string, inboxColl string, retryAttempts int) *Config {
	config := Config{
		MongoDBURL:    mongoDBURL,
		KafkaURL:      kafkaURL,
		OutboxColl:    outboxColl,
		InboxColl:     inboxColl,
		RetryAttempts: retryAttempts,
	}

	if config.MongoDBURL == "" {
		config.MongoDBURL = defaultMongoDBURL
	}
	if config.KafkaURL == "" {
		config.KafkaURL = defaultKafkaURL
	}
	if config.OutboxColl == "" {
		config.OutboxColl = defaultOutboxColl
	}
	if config.InboxColl == "" {
		config.InboxColl = defaultInboxColl
	}
	if config.RetryAttempts == 0 {
		config.RetryAttempts = defaultRetryAttempts
	}

	return &config
}
