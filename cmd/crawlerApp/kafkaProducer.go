package main

import (
	"context"
	"cpmiFeed/pkg/kafkaConfig"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaProducer interface {
	Start()
	Stop()
}

type DefaultProducer struct {
	writer   *kgo.Client
	stopChan chan struct{}
	app      *App
	mu       sync.Mutex
	started  bool
	total    int
	topic    string
}

func (p *DefaultProducer) Start() {
	p.mu.Lock()
	if p.started {
		p.mu.Unlock()
		return
	}
	p.started = true
	p.stopChan = make(chan struct{})
	p.mu.Unlock()

	for {
		select {
		case <-p.stopChan:
			p.mu.Lock()
			defer p.mu.Unlock()
			p.started = false
			return
		case e := <-p.app.eventsChan:
			messages := make([]*kgo.Record, len(e))
			for i, event := range e {
				data, err := json.Marshal(event)
				if err != nil {
					slog.Error("Error marshalling data", "error", err)
				}
				messages[i] = &kgo.Record{
					Value: data,
					Topic: p.topic,
				}
			}
			p.total += len(messages)

			results := p.writer.ProduceSync(context.Background(), messages...)
			for _, result := range results {
				if err := result.Err; err != nil {
					slog.Error("record had a produce error while synchronously producing", "error", err)
				}
			}
		}
	}
}

func (p *DefaultProducer) Stop() {
	if !p.started {
		return
	}
	p.stopChan <- struct{}{}
	close(p.stopChan)
}

func NewKafkaProducer(cfg *kafkaConfig.Config, app *App) KafkaProducer {

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumeTopics(cfg.EventsTopic),
	)
	if err != nil {
		panic(err)
	}
	//defer cl.Close()

	// Connect to Kafka to discover topics
	// conn, err := kafka.Dial("tcp", cfg.Brokers[0])
	// if err != nil {
	// 	slog.Error("Failed to connect to Kafka", "error", err)
	// }
	// defer conn.Close()

	// br, err := conn.Brokers()
	// if err != nil {
	// 	slog.Error("Failed to get the broker metadata", "error", err)
	// }

	// for _, b := range br {
	// 	slog.Info("Broker", "Host", b.Host)
	// }

	return &DefaultProducer{
		writer: cl,
		app:    app,
		topic:  cfg.EventsTopic,
	}
}
