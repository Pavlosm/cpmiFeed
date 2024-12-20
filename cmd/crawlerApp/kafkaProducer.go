package main

import (
	"context"
	"cpmiFeed/pkg/kafka"
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

func NewKafkaProducer(cfg *kafka.Config, app *App) KafkaProducer {

	cl, err := kafka.NewKafkaProducerClient(cfg)

	if err != nil {
		panic(err)
	}

	return &DefaultProducer{
		writer: cl,
		app:    app,
		topic:  cfg.EventsTopic,
	}
}
