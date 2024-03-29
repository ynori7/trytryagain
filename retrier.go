package trytryagain

import (
	"context"
	"fmt"
	"time"
)

var (
	// ErrContextCanceled is returned if the context is cancelled or the deadline is reached
	ErrContextCanceled = fmt.Errorf("context cancelled")
	// ErrRequestNotRetriable is returned if the request fails and the error is not retriable (a 400, for example)
	ErrRequestNotRetriable = fmt.Errorf("request is not retriable")
	// ErrNotSuccessful is returned if the number of retries is exhausted without any successes
	ErrNotSuccessful = fmt.Errorf("request not successful")
)

type (
	// OnErrorFunc is a callback function which is called every time the action returns an error
	OnErrorFunc func(err error)
	// ActionFunc is the function which is called and retried in case of a failure
	ActionFunc func() (err error, retriable bool)
	// ErrorHandler is a function which is called before returning an error
	ErrorHandler func(retrierErr, actionErr error) error
)

// Retrier performs a specified action with retry logic based on its configuration
type Retrier struct {
	maxAttempts  uint
	backoff      BackoffFunc
	onError      OnErrorFunc
	ignoreCtx    bool
	errorHandler ErrorHandler
}

// NewRetrier returns a new Retrier with the specified options
func NewRetrier(options ...RetrierOption) *Retrier {
	t := defaultRetrier()

	for _, opt := range options {
		opt(t)
	}

	return t
}

// Do performs the specified action and retries with backoff in case of failures until the request either succeeds
// or the maximum number of retries has been reached.
func (t *Retrier) Do(ctx context.Context, action ActionFunc) error {
	var actionErr error
	for attempts := uint(0); attempts < t.maxAttempts; attempts++ {
		//sleep for a bit to avoid bombarding the requested resource. The backoff func should return 0 for the first attempt
		time.Sleep(t.backoff(ctx, attempts))

		//check if the context was cancelled
		if IsContextDone(ctx) && !t.ignoreCtx {
			return ErrContextCanceled
		}

		err, retriable := action()
		if err == nil {
			return nil //success
		}

		t.onError(err) //allow the user to handle/log the error

		//it can happen that the context is canceled during the request
		if IsCanceledContextError(err) && !t.ignoreCtx {
			return t.errorHandler(ErrContextCanceled, err)
		}

		if !retriable {
			return t.errorHandler(ErrRequestNotRetriable, err)
		}

		actionErr = err
	}

	return t.errorHandler(ErrNotSuccessful, actionErr)
}
