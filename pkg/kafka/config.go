package kafka

import (
	"os"
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
)

type Config struct {
	Brokers         []string
	EventsTopic     string
	ConsumerGroupId string
	Username        string
	Password        string
	Protocol        string
}

func NewConfig() *Config {
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	eventsTopic := os.Getenv("KAFKA_EVENTS_TOPIC")
	consumerGroupId := os.Getenv("KAFKA_CONSUMER_GROUP_ID")
	return &Config{
		Brokers:         brokers,
		EventsTopic:     eventsTopic,
		ConsumerGroupId: consumerGroupId,
		Username:        os.Getenv("KAFKA_USERNAME"),
		Password:        os.Getenv("KAFKA_PASSWORD"),
		Protocol:        os.Getenv("KAFKA_PROTOCOL"),
	}
}

func NewKafkaConsumerClient(cfg *Config) (*kgo.Client, error) {
	opts := NewKafkaOptions(cfg)
	opts = append(opts, kgo.ConsumerGroup(cfg.ConsumerGroupId))
	return kgo.NewClient(opts...)
}

func NewKafkaProducerClient(cfg *Config) (*kgo.Client, error) {
	opts := NewKafkaOptions(cfg)
	return kgo.NewClient(opts...)
}

func NewKafkaOptions(cfg *Config) []kgo.Opt {
	opts := make([]kgo.Opt, 0)
	opts = append(opts, kgo.SeedBrokers(cfg.Brokers...))
	opts = append(opts, kgo.ConsumeTopics(cfg.EventsTopic))

	if cfg.Username != "" && cfg.Password != "" {
		opts = append(opts, kgo.SASL(
			plain.Auth{
				User: cfg.Username,
				Pass: cfg.Password,
			}.AsMechanism(),
		))
	}

	if cfg.Protocol == "TLS" {
		opts = append(opts, kgo.DialTLS())
	}
	return opts
}
