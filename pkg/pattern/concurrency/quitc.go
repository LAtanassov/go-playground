package concurrency

import (
	"fmt"
	"time"
)

// Aggregator is an example for a quit channel to stop long running task.
type Aggregator struct {
}

// Start has a quit channel to interrupt the endless loop.
func (a *Aggregator) Start(quitc <-chan bool) {
	for {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("aggregating...")
		case <-quitc:
			return
		}
	}
}

// StartWithResource start a long runnning task with resources which need to be cleanup before it can stop.
func (a *Aggregator) StartWithResource(quitc chan string) {
	for {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("aggregating...")
		case <-quitc:
			a.cleanup()
			quitc <- "See you!"
			return
		}

	}
}

func (a *Aggregator) cleanup() {
	fmt.Println("cleaning...")
}
