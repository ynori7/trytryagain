package trytryagain

import (
	"context"
	"math"
	"time"
)

// BackoffFunc is the backoff strategy. Given the number of attempts already made, it returns the wait time
type BackoffFunc func(ctx context.Context, attempts uint) time.Duration

// exponential backoff
func exponentialBackoff(_ context.Context, attempts uint) time.Duration {
	if attempts == 0 {
		return time.Duration(0)
	}
	return time.Duration(math.Pow(10, float64(attempts))) * time.Millisecond
}
