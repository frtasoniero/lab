package utils

import (
	"context"
	"time"
)

// CreateContextServerWithTimeout creates a new context with a 5-second timeout.
// It is typically used to control the duration of operations like database inserts,
// ensuring they don't hang indefinitely.
//
// The function also starts a goroutine to call cancel() once the context is done,
// which helps in cleaning up resources associated with the context.
//
// Returns:
//   - A context.Context that will be automatically canceled after 5 seconds.
func CreateContextServerWithTimeout() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	go func() {
		<-ctx.Done()
		cancel()
	}()
	return ctx
}
