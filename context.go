package trytryagain

import (
	"context"
	"strings"
)

// IsCanceledContextError returns true if the error is some variant of "context canceled" or "context deadline exceeded"
func IsCanceledContextError(err error) bool {
	errLower := strings.ToLower(err.Error())
	return strings.Contains(errLower, context.Canceled.Error()) || strings.Contains(errLower, context.DeadlineExceeded.Error())
}

// IsContextDone returns true if the context is done
func IsContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
