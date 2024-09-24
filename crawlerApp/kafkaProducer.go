package main

import (
	"context"
	"cpmiFeed/kafkaConfig"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
	Start()
	Stop()
}

type DefaultProducer struct {
	writer   *kafka.Writer
	stopChan chan struct{}
	app      *App
	mu       sync.Mutex
	started  bool
	total    int
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
			messages := make([]kafka.Message, len(e))
			for i, event := range e {
				data, err := json.Marshal(event)
				if err != nil {
					slog.Error("Error marshalling data", "error", err)
				}
				messages[i] = kafka.Message{
					Value: data,
				}
			}
			p.total += len(messages)

			err := p.writer.WriteMessages(context.Background(), messages...)
			if err != nil {
				slog.Error("Error writing messages", "error", err)
			}
			slog.Info("Written messages", "number", len(messages), "total", p.total)
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
	return &DefaultProducer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: cfg.Brokers,
			Topic:   cfg.EventsTopic,
		}),
		app: app,
	}
}
