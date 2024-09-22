package main

import (
	"context"
	"cpmiFeed/rawEventModels"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer interface {
	Start(onNewMessage func(events []rawEventModels.Event) error)
	Stop()
}

type DefaultConsumer struct {
	reader   *kafka.Reader
	stopChan chan struct{}
	app      *App
	mu       sync.Mutex
	started  bool
}

func NewDefaultConsumer(brokers []string, topic string, app *App) KafkaConsumer {
	return &DefaultConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: "cpmiEventsConsumer",
		}),
		app: app,
	}
}

func (c *DefaultConsumer) Start(onNewMessage func(events []rawEventModels.Event) error) {
	c.mu.Lock()
	if c.started {
		c.mu.Unlock()
		return
	}
	c.started = true
	c.stopChan = make(chan struct{})
	c.mu.Unlock()

	for {
		select {
		case <-c.stopChan:
			c.mu.Lock()
			defer c.mu.Unlock()
			c.started = false
			return
		default:
			m, err := c.reader.FetchMessage(context.Background())
			if err != nil {
				slog.Error("Error reading message", "error", err)
				continue
			}

			var event rawEventModels.Event
			err = json.Unmarshal(m.Value, &event)
			if err != nil {
				slog.Error("Error unmarshalling message", "error", err)
				continue
			}

			// todo on commit
			err = onNewMessage([]rawEventModels.Event{event})
			if err != nil {
				slog.Error("Error processing message", "error", err)
				continue
			}
			c.reader.CommitMessages(context.Background(), m)

			c.app.eventsChan <- []rawEventModels.Event{event}
		}
	}
}

func (c *DefaultConsumer) Stop() {
	if !c.started {
		return
	}
	c.stopChan <- struct{}{}
	err := c.reader.Close()
	if err != nil {
		slog.Error("Error closing reader", "error", err)
	}

}
