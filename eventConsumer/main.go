package main

import (
	"cpmiFeed/common"
	"cpmiFeed/db"
	"cpmiFeed/kafkaConfig"
	"log/slog"
	"os"
	"sync"
)

func main() {
	app := App{
		wg:         &sync.WaitGroup{},
		eventsChan: make(chan []common.Event),
	}

	stop := make(chan os.Signal, 1)

	repos := db.NewRepositories()
	defer repos.Close()

	cfg := kafkaConfig.NewConfig()
	consumer := NewDefaultConsumer(cfg, &app)

	go consumer.Start(repos.Event.Save)

	m := 0
	go func() {
		for {
			events := <-app.eventsChan
			m += len(events)
			slog.Info("Received events", "messageNo", m, "events", events)
		}
	}()

	<-stop

	consumer.Stop()
	app.wg.Wait()
}

// mongodb+srv://admin:FmGXU6j1kPvT6ovb@cpmi-crawler-cluster.wlsq1.mongodb.net/cpmiFeed?retryWrites=true&w=majority
// mongodb+srv://admin:FmGXU6j1kPvT6ovb@cpmi-crawler-cluster.wlsq1.mongodb.net/?retryWrites=true&w=majority&appName=cpmi-crawler-cluster
