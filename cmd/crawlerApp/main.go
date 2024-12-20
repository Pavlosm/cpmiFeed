package main

import (
	"cpmiFeed/pkg/common"
	"cpmiFeed/pkg/kafka"
	"os"
	"sync"
)

func main() {
	app := App{
		wg:         &sync.WaitGroup{},
		eventsChan: make(chan []common.Event, 5000),
	}

	stop := make(chan os.Signal, 1)

	cfg := kafka.NewConfig()
	producer := NewKafkaProducer(cfg, &app)
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
