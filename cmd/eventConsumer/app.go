package main

import (
	"cpmiFeed/pkg/common"
	"sync"
)

type App struct {
	wg         *sync.WaitGroup
	eventsChan chan []common.Event
}
