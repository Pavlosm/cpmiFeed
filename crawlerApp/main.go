package main

import (
	"cpmiFeed/rawEventModels"
	"os"
	"sync"
)

const (
	kafkaURL = "localhost:29092"
	topic    = "cpmiEvents"
)

func main() {
	app := App{
		wg:         &sync.WaitGroup{},
		eventsChan: make(chan []rawEventModels.Event, 5000),
	}

	stop := make(chan os.Signal, 1)

	producer := NewKafkaProducer([]string{kafkaURL}, topic, app.eventsChan, &app)
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
