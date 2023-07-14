package trytryagain

var (
	_defaultMaxAttmpts = uint(3)
	_defaultBackoff    = exponentialBackoff
	_defaultOnError    = func(err error) {}
	_defaultIgnoreCtx  = false
)

// RetrierOption is a callback for specifying configuration options for a Retrier
type RetrierOption func(t *Retrier)

func defaultRetrier() *Retrier {
	return &Retrier{
		maxAttempts: _defaultMaxAttmpts,
		backoff:     _defaultBackoff,
		onError:     _defaultOnError,
		ignoreCtx:   _defaultIgnoreCtx,
	}
}

// WithMaxAttempts is an option to specify the maximum number of retries which the Retrier should make
func WithMaxAttempts(maxAttempts uint) RetrierOption {
	return func(r *Retrier) {
		r.maxAttempts = maxAttempts
	}
}

// WithBackoff allows you to specify the backoff strategy, for example an exponential backoff
func WithBackoff(backoff BackoffFunc) RetrierOption {
	return func(r *Retrier) {
		r.backoff = backoff
	}
}

// WithOnError is an option to specify the OnError callback which is called for each failed attempt
func WithOnError(onError OnErrorFunc) RetrierOption {
	return func(r *Retrier) {
		r.onError = onError
	}
}

// WithIgnoreCtx is an option to specify whether context canncelation/deadline should be ignored (will retry anyway if true)
func WithIgnoreCtx(ignoreCtx bool) RetrierOption {
	return func(r *Retrier) {
		r.ignoreCtx = ignoreCtx
	}
}
