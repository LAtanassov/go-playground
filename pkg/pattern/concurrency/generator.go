package concurrency

import (
	"math/rand"
	"time"
)

// Generater specifies interface of a function that returns a channel
type Generater interface {
	Generate() <-chan struct{}
}

// Noise generates randomness.
type Noise struct {
}

// Generate random integers through a channel
func (n *Noise) Generate() <-chan int {
	c := make(chan int)
	go func() {
		for i := 0; ; i++ {
			c <- rand.Intn(1e3)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c
}
