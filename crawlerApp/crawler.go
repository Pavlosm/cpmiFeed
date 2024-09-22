package main

import (
	"cpmiFeed/rawEventModels"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"
)

type Crawler interface {
	Start() error
	ShutDown() error
}

type stubCrawler struct {
	app         *App
	started     bool
	mu          sync.Mutex
	numOfCrawls int
	stopChan    chan struct{}
	stoppedChan chan struct{}
}

func NewStubCrawler(app *App) Crawler {
	return &stubCrawler{
		app: app,
		mu:  sync.Mutex{},
	}
}

func (c *stubCrawler) Start() error {
	c.mu.Lock()
	if c.started {
		c.mu.Unlock()
		return nil
	}

	c.started = true
	c.stopChan = make(chan struct{})
	c.stoppedChan = make(chan struct{})

	c.mu.Unlock()

	for {
		select {
		case <-c.stopChan:
			c.mu.Lock()
			defer c.mu.Unlock()
			c.started = false
			c.stoppedChan <- struct{}{}
			return nil
		default:
			c.crawl()
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (c *stubCrawler) ShutDown() error {
	if !c.started {
		return nil
	}
	c.stopChan <- struct{}{}
	<-c.stoppedChan
	close(c.stopChan)
	close(c.stoppedChan)
	return nil
}

func (c *stubCrawler) crawl() {
	if !c.started {
		return
	}
	c.app.wg.Add(1)
	defer c.app.wg.Done()

	c.numOfCrawls++

	r := rand.Intn(10)
	events := make([]rawEventModels.Event, r)
	t := make([]string, r)
	tags := []string{"Education", "Conference", "AI", "C#", "Golang"}

	for i := 0; i < r; i++ {
		if i%2 == 0 {
			rt1 := rand.Intn(len(tags))
			rt2 := rand.Intn(len(tags))
			t = append(t, tags[rt1], tags[rt2])
		} else if i%3 == 0 {
			rt1 := rand.Intn(len(tags))
			rt2 := rand.Intn(len(tags))
			rt3 := rand.Intn(len(tags))
			t = append(t, tags[rt1], tags[rt2], tags[rt3])
		} else {
			rt := rand.Intn(len(tags))
			t = append(t, tags[rt])
		}
		events = append(events, rawEventModels.Event{
			Data:        fmt.Sprintf("Event %d", i),
			URL:         fmt.Sprintf("https://example.com/%d", i),
			Description: fmt.Sprintf("Description %d", i),
			Tags:        t,
			Timestamp:   time.Now(),
		})
	}

	slog.Info("Crawled events", "number", len(events))
	c.app.eventsChan <- events
}
