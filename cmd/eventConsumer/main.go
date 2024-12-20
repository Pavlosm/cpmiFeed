package main

import (
	"cpmiFeed/pkg/common"
	"cpmiFeed/pkg/db"
	"cpmiFeed/pkg/kafka"
	"log/slog"
	"os"
	"sync"
	"time"
)

func main() {
	app := App{
		wg:         &sync.WaitGroup{},
		eventsChan: make(chan []common.Event, 1000),
	}

	stop := make(chan os.Signal, 1)

	repos := db.NewRepositories()
	defer repos.Close()

	cfg := kafka.NewConfig()
	consumer := NewDefaultConsumer(cfg, &app)

	go consumer.Start(repos.Event.Save)

	m := 0
	filterer, err := NewUserEventFilterer(repos)
	if err != nil {
		panic("could not get user filters")
	}

	go func() {
		for {
			events := <-app.eventsChan
			time.Sleep(time.Millisecond * 10)
			go filterer.HandleSafely(events)
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
