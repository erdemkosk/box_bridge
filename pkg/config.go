package pkg

type BoxBridgeConfig struct {
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

type ConfigBuilder struct {
	config *BoxBridgeConfig
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{config: &BoxBridgeConfig{}}
}

func (cb *ConfigBuilder) WithMongoDBURL(url string) *ConfigBuilder {
	cb.config.MongoDBURL = url
	return cb
}

func (cb *ConfigBuilder) WithKafkaURL(url string) *ConfigBuilder {
	cb.config.KafkaURL = url
	return cb
}

func (cb *ConfigBuilder) WithOutboxCollection(collection string) *ConfigBuilder {
	cb.config.OutboxColl = collection
	return cb
}

func (cb *ConfigBuilder) WithInboxCollection(collection string) *ConfigBuilder {
	cb.config.InboxColl = collection
	return cb
}

func (cb *ConfigBuilder) WithRetryAttempts(attempts int) *ConfigBuilder {
	cb.config.RetryAttempts = attempts
	return cb
}

func (cb *ConfigBuilder) Build() *BoxBridgeConfig {
	if cb.config.MongoDBURL == "" {
		cb.config.MongoDBURL = defaultMongoDBURL
	}
	if cb.config.KafkaURL == "" {
		cb.config.KafkaURL = defaultKafkaURL
	}
	if cb.config.OutboxColl == "" {
		cb.config.OutboxColl = defaultOutboxColl
	}
	if cb.config.InboxColl == "" {
		cb.config.InboxColl = defaultInboxColl
	}
	if cb.config.RetryAttempts == 0 {
		cb.config.RetryAttempts = defaultRetryAttempts
	}

	return cb.config
}
