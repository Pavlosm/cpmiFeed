package main

import (
	"cpmiFeed/rawEventModels"
	"sync"
)

type App struct {
	wg         *sync.WaitGroup
	eventsChan chan []rawEventModels.Event
}
