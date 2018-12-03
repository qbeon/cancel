package cancel_test

import (
	"context"
	"test/cancel"
	"testing"
)

var gtk cancel.Token
var gctx context.Context
var gb bool
var ge error

func takeCancelable(token cancel.Token) {
	gtk = token
}

func takeCancelableCtx(ctx context.Context) {
	gctx = ctx
}

func BenchmarkCreation(b *testing.B) {
	gen := cancel.NewGenerator()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tk := gen.New()
		tk.Cancel()
	}
}

func BenchmarkCreationCtx(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, cancel := context.WithCancel(context.Background())
		cancel()
	}
}

func BenchmarkCopy(b *testing.B) {
	gen := cancel.NewGenerator()
	tk := gen.New()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		gtk = tk
	}
	tk.Cancel()
}

func BenchmarkCopyCtx(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		gctx = ctx
	}
	cancel()
}

func BenchmarkIsCancelled(b *testing.B) {
	gen := cancel.NewGenerator()
	tk := gen.New()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		gb = tk.IsCanceled()
	}
	tk.Cancel()
}

func BenchmarkIsCancelledCtx(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ge = ctx.Err()
	}
	cancel()
}

func BenchmarkChan(b *testing.B) {
	gen := cancel.NewGenerator()
	tk := gen.New()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		select {
		case <-tk.Canceled():
		default:
		}
	}
	tk.Cancel()
}

func BenchmarkChanCtx(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		select {
		case <-ctx.Done():
		default:
		}
	}
	cancel()
}
