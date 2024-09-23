package main

import (
	"cpmiFeed/common"
	"os"
	"sync"
)

const (
	topic = "cpmiEvents"
)

func main() {
	app := App{
		wg:         &sync.WaitGroup{},
		eventsChan: make(chan []common.Event, 5000),
	}

	stop := make(chan os.Signal, 1)

	producer := NewKafkaProducer([]string{"localhost:29092", "localhost:29093", "localhost:29094"}, topic, app.eventsChan, &app)
	defer producer.Stop()
	go producer.Start()

	crawler := NewStubCrawler(&app)
	defer crawler.ShutDown()
	go crawler.Start()

	<-stop

	crawler.ShutDown()
	producer.Stop()
	app.wg.Wait()
}
