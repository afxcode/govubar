package main

import (
	"context"
	"sync"

	"govubar/darwin/menuet"
)

func getDisplayContext() (wg *sync.WaitGroup, ctx context.Context) {
	return menuet.App().GracefulShutdownHandles()
}

func runDisplay() {
	menuet.App().RunApplication()
}

func setDisplayBarGraph(audioLevel int) {
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: generateBarGraph(audioLevel),
	})
}
