package main

import (
	"log"
	"sync"
	"time"

	"github.com/qbeon/cancel"
)

// Cancelable starts counting indefinitely and returns the number it stopped at
// when the cancelation token is canceled by the caller
func Cancelable(c cancel.Token) uint64 {
	for i := uint64(0); ; i++ {
		// Check whether the cancelation token was canceled
		if c.IsCanceled() {
			return i
		}
	}
}

func main() {
	// Get a new cancelation token from the global pool
	cancelationToken := cancel.New()

	// Cancel the token eventually to return it to the pool
	defer cancelationToken.Cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		i := Cancelable(cancelationToken)
		log.Print("stopped at: ", i)
		wg.Done()
	}()

	// Wait a little and cancel the token
	time.Sleep(500 * time.Millisecond)
	cancelationToken.Cancel()

	// Wait until the counter is stopped
	wg.Wait()
}
