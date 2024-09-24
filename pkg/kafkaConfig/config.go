package kafkaConfig

import (
	"os"
	"strings"
)

type Config struct {
	Brokers         []string
	EventsTopic     string
	ConsumerGroupId string
}

func NewConfig() *Config {
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	eventsTopic := os.Getenv("KAFKA_EVENTS_TOPIC")
	consumerGroupId := os.Getenv("KAFKA_CONSUMER_GROUP_ID")
	return &Config{
		Brokers:         brokers,
		EventsTopic:     eventsTopic,
		ConsumerGroupId: consumerGroupId,
	}
}