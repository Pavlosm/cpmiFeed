package main

import (
	"context"
	"cpmiFeed/pkg/common"
	"cpmiFeed/pkg/db"
	"cpmiFeed/pkg/kafkaConfig"
	"log/slog"
	"os"
	"strings"
	"sync"
)

func main() {
	app := App{
		wg:         &sync.WaitGroup{},
		eventsChan: make(chan []common.Event, 1000),
	}

	stop := make(chan os.Signal, 1)

	repos := db.NewRepositories()
	defer repos.Close()

	cfg := kafkaConfig.NewConfig()
	consumer := NewDefaultConsumer(cfg, &app)

	go consumer.Start(repos.Event.Save)

	m := 0
	userFilters, err := repos.UserFilter.GetAll(context.TODO())
	filterer := ConcreteUserEventFilterer{
		repos:       repos,
		userFIlters: userFilters,
	}

	if err != nil {
		panic("could not get user filters")
	}

	go func() {
		for {
			events := <-app.eventsChan
			go filterer.Handle(events)
			m += len(events)
			slog.Info("Received events", "messageNo", m, "events", events)
			if err != nil {
				continue
			}
		}
	}()

	<-stop

	consumer.Stop()
	app.wg.Wait()
}

type UserEventFilterer interface {
	Init() error
	Handle(events []common.Event)
}

type ConcreteUserEventFilterer struct {
	repos       *db.Repositories
	userFIlters map[string][]common.UserEventFilter
}

func (c *ConcreteUserEventFilterer) Init() error {
	filters, err := c.repos.UserFilter.GetAll(context.TODO())
	if err != nil {
		return err
	}
	c.userFIlters = filters
	return nil
}

func (c *ConcreteUserEventFilterer) Handle(events []common.Event) error {
	for _, e := range events {
		for id, fs := range c.userFIlters {
			if anyFilterApplies(e, fs) {
				err := c.repos.UserEvents.UpsertUserEvents(context.TODO(), id, e)
				if err != nil {
					slog.Error("Could not upsert event", err)
				} else {
					slog.Info("added new user event")
				}
			}
		}
	}
	return nil
}

func anyFilterApplies(e common.Event, fs []common.UserEventFilter) bool {
	for _, f := range fs {
		for _, t := range f.Tags {
			for _, tt := range e.Tags {
				if t == tt {
					return true
				}
			}
		}

		for _, tk := range f.Tokens {
			if strings.Contains(e.Data, tk) || strings.Contains(e.Description, tk) {
				return true
			}
		}
	}
	slog.Info("event skipped")
	return false
}

type NonInitError struct {
	message string
}

func (e *NonInitError) Error() string {
	return e.message
}
