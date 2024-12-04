package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func getDisplayContext() (wg *sync.WaitGroup, ctx context.Context) {
	return &sync.WaitGroup{}, context.Background()
}

func runDisplay() {
	for {
		time.Sleep(time.Second)
	}
}

func setDisplayBarGraph(audioLevel int) {
	fmt.Println(generateBarGraph(audioLevel))
}
