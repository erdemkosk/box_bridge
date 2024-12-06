package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/erdemkosk/box_bridge/pkg"
)

func GetHeaderValue(headers []kafka.Header, key string) (string, bool) {
	for _, header := range headers {
		if header.Key == key {
			return string(header.Value), true
		}
	}
	return "", false
}

func ConvertKafkaHeadersToWrapperHeaders(kafkaHeaders []pkg.KafkaHeader) []kafka.Header {
	headers := make([]kafka.Header, len(kafkaHeaders))
	for i, kh := range kafkaHeaders {
		headers[i] = kafka.Header{
			Key:   kh.Key,
			Value: kh.Value,
		}
	}
	return headers
}

func defaultInt(value, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}

func defaultString(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
