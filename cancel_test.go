package cancel_test

import (
	"sync"
	"test/cancel"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newToken creates a new generator and generates a new token
func newToken() cancel.Token {
	origin := cancel.NewGenerator()
	return origin.New()
}

// TestChannel tests the notification channel
func TestChannel2(t *testing.T) {
	token := newToken()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		c := make(chan struct{})
		select {
		case <-c:
			// this will never trigger
			t.Fail()
		case status := <-token.Canceled():
			assert.Equal(t, struct{}{}, status)
		}
	}()

	time.Sleep(50 * time.Millisecond)
	require.True(t, token.Cancel())

	wg.Wait()
}

// TestCancel tests cancelation
func TestCancel(t *testing.T) {
	token := newToken()
	require.True(t, token.Cancel())
}

// TestIsCancelled tests Token.IsCanceled
func TestIsCancelled(t *testing.T) {
	token := newToken()

	require.False(t, token.IsCanceled())
	require.False(t, token.IsCanceled())

	require.True(t, token.Cancel())

	require.True(t, token.IsCanceled())
	require.True(t, token.IsCanceled())
}

// TestRepeatedCancel tests calling Token.Cancel multiple times in a row
func TestRepeatedCancel(t *testing.T) {
	token := newToken()

	require.True(t, token.Cancel())
	require.False(t, token.Cancel())
	require.False(t, token.Cancel())
}

// TestChannelAfterCancel tests reading the channel after closure
func TestChannelAfterCancel(t *testing.T) {
	token := newToken()

	token.Cancel()

	<-token.Canceled()
	<-token.Canceled()
	<-token.Canceled()
}

func TestC(t *testing.T) {
	origin := cancel.NewGenerator()

	tk := origin.New()
	tk.Cancel()

	tk = origin.New()
	tk.Cancel()
}
