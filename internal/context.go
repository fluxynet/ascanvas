package internal

import (
	"context"
	"io"
)

// IsContextDone is a helper to check if a context.Context is done
func IsContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// Closed is differable closure of Closable items with nil checking
func Closed(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
}
