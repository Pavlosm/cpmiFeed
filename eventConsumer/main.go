package main

import (
	"cpmiFeed/db"
	"cpmiFeed/rawEventModels"
	"log/slog"
	"os"
	"sync"
)

func main() {
	app := App{
		wg:         &sync.WaitGroup{},
		eventsChan: make(chan []rawEventModels.Event),
	}

	stop := make(chan os.Signal, 1)

	repo, err := db.NewMongoEventRepository("mongodb://root:example@localhost:27017", "cpmiFeed")
	if err != nil {
		os.Exit(1)
	}
	defer repo.Close()
	consumer := NewDefaultConsumer([]string{"localhost:29092", "localhost:29093", "localhost:29094"}, "cpmiEvents", &app)
	go consumer.Start(repo.Save)

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
