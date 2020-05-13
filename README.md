# TryTryAgain [![GoDoc](https://godoc.org/github.com/ynori7/trytryagain?status.png)](https://godoc.org/github.com/ynori7/trytryagain) [![Build Status](https://travis-ci.org/ynori7/trytryagain.svg?branch=master)](https://travis-ci.com/github/ynori7/trytryagain) [![Coverage Status](https://coveralls.io/repos/github/ynori7/trytryagain/badge.svg?branch=master)](https://coveralls.io/github/ynori7/trytryagain?branch=master) [![Go Report Card](https://goreportcard.com/badge/ynori7/trytryagain)](https://goreportcard.com/report/github.com/ynori7/trytryagain)
This library provides a simple utility for performing an action with retries. TryTryAgain is thread-safe and takes 
care of handling backoff and expired contexts. The library provides configuration for the backoff strategy, retry 
counts, and a callback for errors.

# How it works
A retrier is created with the specified options for maxAttempts, backoffFunc, and onErrorFunc. Then you simply call `Do`, 
providing a function which wraps the action to be performed. This function should return an error and a boolean to indicate 
whether the error is retriable or not. 

**Defaults:** 

- The default max attempts is 3
- The default backoff strategy is exponential 
- The default onError callback does nothing

# Usage

Here is a trivial example:

```go
r := NewRetrier(
    WithMaxAttempts(4),
    WithOnError(func(err error) { fmt.Println(err.Error()) }),
    WithBackoff(exponentialBackoff),
)

err := r.Do(ctx, func() (error, bool) {
    return fmt.Errorf("something went wrong"), true
})
```

More detailed examples can be found in [examples](./examples)
