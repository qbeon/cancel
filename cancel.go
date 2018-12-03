package cancel

import (
	"sync"
	"sync/atomic"
)

/* GLOBAL GENERTOR */

var globalGenerator *Generator

func init() {
	// Initialize the global token pool
	globalGenerator = NewGenerator()
}

// New issues a new cancelation token instance that must be canceled through
// token.Cancel() when it's no longer needed to avoid leaking memory
func New() Token {
	return globalGenerator.New()
}

// Generator represents a cancelation token generator
type Generator struct {
	// pool holds pooled token references
	pool sync.Pool

	// identCounter stores the last assigned identifier
	identCounter uint64
}

/* GENERATOR */

// NewGenerator creates a new generator instance
func NewGenerator() *Generator {
	newGenerator := &Generator{
		pool:         sync.Pool{},
		identCounter: 0,
	}

	newGenerator.pool.New = func() interface{} {
		return &token{
			origin:  newGenerator,
			channel: make(chan struct{}, 1),
		}
	}

	return newGenerator
}

// New issues a new cancelation token instance that must be canceled through
// token.Cancel() when it's no longer needed to avoid leaking memory
func (gen *Generator) New() Token {
	tokenIdent := atomic.AddUint64(&gen.identCounter, 1)

	// Get a token from the pool
	token := gen.pool.Get().(*token)

	// Assign a new unique identifier
	atomic.StoreUint64(&token.ident, tokenIdent)

	// Drain the channel
	select {
	case <-token.channel:
	default:
	}

	// Wrap the token
	return Token{
		token: token,
		ident: tokenIdent,
	}
}

/* TOKEN */

// Token represents a thread-safe stateless cancelation token wrapper, it can
// safely be copied shallowly
type Token struct {
	ident uint64
	token *token
}

// isClosed determines whether the token was closed by comparing the identifiers
// of the shared token instance and the token wrapper, it those don't match then
// the wrapper doesn't reference its original token any longer which means the
// actual token was already closed
func (wrp Token) isClosed() bool {
	return wrp.ident != atomic.LoadUint64(&wrp.token.ident)
}

// Cancel returns true if the token was canceled, otherwise returns false
// indicating that the token was already canceled or closed
func (wrp Token) Cancel() bool {
	if wrp.isClosed() {
		return false
	}

	// Reset the identifier to zero
	atomic.StoreUint64(&wrp.token.ident, 0)

	// Notify channel listeners about the closure
	wrp.token.channel <- struct{}{}

	// Return the token to the pool
	wrp.token.origin.pool.Put(wrp.token)
	return true
}

// IsCanceled returns true if the token is already canceled, otherwise returns
// false
func (wrp Token) IsCanceled() bool {
	return wrp.isClosed()
}

// Canceled returns a read-only channel that's triggered when the token is
// canceled. If the token was already canceled a dummy channel is allocated
func (wrp Token) Canceled() <-chan struct{} {
	if wrp.isClosed() {
		c := make(chan struct{}, 1)
		c <- struct{}{}
		return c
	}
	return wrp.token.channel
}

// token represents the internal shared token
type token struct {
	origin  *Generator
	ident   uint64
	channel chan struct{}
}
