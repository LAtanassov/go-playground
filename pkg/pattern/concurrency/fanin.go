package concurrency

import (
	"time"
)

// FanIner specifies the interface of the fan in pattern.
type FanIner interface {
	FanIn(a, b <-chan int) <-chan int
}

// Pipeline is reference implementation of fan in patterns.
type Pipeline struct {
}

// FanIn synchronzing two channels into one channel.
func (s *Pipeline) FanIn(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			c <- <-a
		}
	}()
	go func() {
		for {
			c <- <-b
		}
	}()
	return c
}

// Highway is just another example of the fan in pattern.
type Highway struct {
}

// FanIn synchronzing two channels into one channel with a select.
func (h *Highway) FanIn(a, b <-chan int) <-chan int {
	c := make(chan int)
	// 5 seconds from now on
	timeout := time.After(5 * time.Second)
	go func() {
		for {
			select {
			case s := <-a:
				c <- s
			case s := <-b:
				c <- s
			// 1 second when select called
			case <-time.After(1 * time.Second):
				c <- 0
			case <-timeout:
				c <- 1
			}
		}
	}()
	return c
}
