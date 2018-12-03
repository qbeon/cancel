package main

import (
	"log"
	"sync"
	"time"

	"github.com/qbeon/cancel"
)

// CancelableCounter starts counting indefinitely and returns the number it
// stopped at when the cancelation token is canceled by the caller
func CancelableCounter(c cancel.Token) uint64 {
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
	wg.Add(2)

	// Start the counter goroutine
	go func() {
		i := CancelableCounter(cancelationToken)
		log.Print("stopped at: ", i)
		wg.Done()
	}()

	// Start a goroutine that's going to wait asynchronously until the
	// cancelation token is canceled
	go func() {
		c := make(chan struct{})
		select {
		case <-c:
			// This will never trigger
		case <-cancelationToken.Canceled():
			// This will trigger when the cancelation token is canceled
			log.Print("canceled!")
			wg.Done()
		}
	}()

	// Wait a little and cancel the token
	time.Sleep(500 * time.Millisecond)
	cancelationToken.Cancel()

	// Wait until the counter is stopped
	wg.Wait()
}
