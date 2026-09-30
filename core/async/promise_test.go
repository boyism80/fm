package async

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestThenRPCCallsOnce(t *testing.T) {
	for i := 0; i < 200; i++ {
		var calls atomic.Int32
		done := make(chan struct{})
		ThenRPC(NewPromise(nil, time.Second), func(context.Context) (int, error) {
			calls.Add(1)
			return 1, nil
		}, func(int) error {
			return nil
		}).Finally(func() {
			close(done)
		})

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatalf("iteration %d: promise did not finish", i)
		}
		time.Sleep(time.Millisecond)
		if n := calls.Load(); n != 1 {
			t.Fatalf("iteration %d: rpc called %d times", i, n)
		}
	}
}
