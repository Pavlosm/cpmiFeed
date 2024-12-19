package main

import (
	"context"
	"cpmiFeed/pkg/common"
	"cpmiFeed/pkg/kafkaConfig"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaConsumer interface {
	Start(onNewMessage func(events []common.Event) error)
	Stop()
}

type DefaultConsumer struct {
	reader   *kgo.Client
	stopChan chan struct{}
	app      *App
	mu       sync.Mutex
	started  bool
	topic    string
}

func NewDefaultConsumer(cfg *kafkaConfig.Config, app *App) KafkaConsumer {
	cl, err := kafkaConfig.NewKafkaConsumerClient(cfg)

	if err != nil {
		panic(err)
	}

	return &DefaultConsumer{
		reader: cl,
		app:    app,
		topic:  cfg.EventsTopic,
	}
}

func (c *DefaultConsumer) Start(onNewMessage func(events []common.Event) error) {
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

			fetches := c.reader.PollRecords(context.Background(), 1)
			if fetches.IsClientClosed() {
				return
			}
			if errs := fetches.Errors(); len(errs) > 0 {
				for _, err := range errs {
					slog.Error("Error fetching messages", "error", err)
				}
				continue
			}
			records := fetches.Records()
			if len(records) == 0 {
				continue
			}
			m := records[0]

			var event common.Event
			err := json.Unmarshal(m.Value, &event)
			event.ID = fmt.Sprintf("%d-%d", m.Partition, m.Offset)
			if err != nil {
				slog.Error("Error unmarshalling message", "error", err)
				continue
			}

			err = onNewMessage([]common.Event{event})
			if err != nil {
				slog.Error("Error processing message", "error", err)
				continue
			}
			c.reader.CommitRecords(context.Background(), fetches.Records()...)

			c.app.eventsChan <- []common.Event{event}
		}
	}
}

func (c *DefaultConsumer) Stop() {
	if !c.started {
		return
	}
	c.stopChan <- struct{}{}
	c.reader.Close()
}
