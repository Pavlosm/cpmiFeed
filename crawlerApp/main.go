package main

import (
	"cpmiFeed/common"
	"cpmiFeed/kafkaConfig"
	"os"
	"sync"
)

func main() {
	app := App{
		wg:         &sync.WaitGroup{},
		eventsChan: make(chan []common.Event, 5000),
	}

	stop := make(chan os.Signal, 1)

	cfg := kafkaConfig.NewConfig()
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
