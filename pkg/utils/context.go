package utils

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func WithContextSigtermCallback(ctx context.Context, f func()) context.Context {
	ctx, cancel := context.WithCancelCause(ctx)
	go func() {
		receivedSignal := make(chan os.Signal, 1)
		signal.Notify(receivedSignal, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(receivedSignal)

		select {
		case <-ctx.Done():
		case <-receivedSignal:
			if f != nil {
				f()
			}
			cancel(fmt.Errorf("received signal: %v", receivedSignal))
		}
	}()

	return ctx
}
