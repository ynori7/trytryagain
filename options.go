package trytryagain

var (
	_defaultMaxAttmpts = uint(3)
	_defaultBackoff    = exponentialBackoff
	_defaultOnError    = func(err error) {}
)

type retrierOption func(t *retrier)

func defaultRetrier() *retrier {
	return &retrier{
		maxAttempts: _defaultMaxAttmpts,
		backoff:     _defaultBackoff,
		onError:     _defaultOnError,
	}
}

// WithMaxAttempts is an option to specify the maximum number of retries which the retrier should make
func WithMaxAttempts(maxAttempts uint) retrierOption {
	return func(r *retrier) {
		r.maxAttempts = maxAttempts
	}
}

// WithBackoff allows you to specify the backoff strategy, for example an exponential backoff
func WithBackoff(backoff BackoffFunc) retrierOption {
	return func(r *retrier) {
		r.backoff = backoff
	}
}

// WithOnError is an option to specify the OnError callback which is called for each failed attempt
func WithOnError(onError OnErrorFunc) retrierOption {
	return func(r *retrier) {
		r.onError = onError
	}
}
