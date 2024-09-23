package main

import (
	"cpmiFeed/common"
	"sync"
)

type App struct {
	wg         *sync.WaitGroup
	eventsChan chan []common.Event
}
