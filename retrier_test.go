package trytryagain

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DoWithRetries(t *testing.T) {
	// given
	var (
		actionCalledCount int
		errCalledCount    int
	)

	resetCounts := func() {
		actionCalledCount = 0
		errCalledCount = 0
	}

	defaultOnError := func(err error) {
		errCalledCount++
		fmt.Println("Got an error: ", err.Error())
	}

	testcases := map[string]struct {
		action              ActionFunc
		onErr               OnErrorFunc
		expectedActionCount int
		expectedErrCount    int
		expectedFinalErr    error
	}{
		"Success after some errors": {
			action: func() (error, bool) {
				actionCalledCount++
				if actionCalledCount > 2 {
					return nil, false
				}
				return fmt.Errorf("something went wrong"), true
			},
			onErr:               defaultOnError,
			expectedActionCount: 3,
			expectedErrCount:    2,
			expectedFinalErr:    nil,
		},
		"Immediate success": {
			action: func() (error, bool) {
				actionCalledCount++
				return nil, false
			},
			onErr:               defaultOnError,
			expectedActionCount: 1,
			expectedErrCount:    0,
			expectedFinalErr:    nil,
		},
		"Non-retriable error": {
			action: func() (error, bool) {
				actionCalledCount++
				return fmt.Errorf("something went wrong"), false
			},
			onErr:               defaultOnError,
			expectedActionCount: 1,
			expectedErrCount:    1,
			expectedFinalErr:    ErrRequestNotRetriable,
		},
		"Failure": {
			action: func() (error, bool) {
				actionCalledCount++
				return fmt.Errorf("something went wrong"), true
			},
			onErr:               defaultOnError,
			expectedActionCount: 4,
			expectedErrCount:    4,
			expectedFinalErr:    ErrNotSuccessful,
		},
	}

	for testcase, testdata := range testcases {
		// when
		tr := NewRetrier(
			WithMaxAttempts(4),
			WithOnError(testdata.onErr),
			WithBackoff(exponentialBackoff),
		)
		err := tr.Do(context.Background(), testdata.action)

		// then
		if testdata.expectedFinalErr == nil {
			require.NoError(t, err, testcase)
		} else {
			require.Error(t, err, testcase)
			assert.EqualError(t, testdata.expectedFinalErr, err.Error())
		}

		assert.Equal(t, testdata.expectedActionCount, actionCalledCount)
		assert.Equal(t, testdata.expectedErrCount, errCalledCount)

		// cleanup
		resetCounts()
	}
}

func Test_DoWithRetries_ContextCanceledAfterFirstAttempt(t *testing.T) {
	// given
	ctx, cancel := context.WithCancel(context.Background())
	actionCalledCount := 0
	action := func() (error, bool) {
		actionCalledCount++
		cancel()
		return fmt.Errorf("something went wrong"), true
	}

	errCalledCount := 0
	onError := func(err error) {
		errCalledCount++
		fmt.Println("Got an error: ", err.Error())
	}

	tr := NewRetrier(
		WithOnError(onError),
	)

	// when
	err := tr.Do(ctx, action)

	// then
	require.Error(t, err)
	assert.EqualError(t, ErrContextCanceled, err.Error())
	assert.Equal(t, 1, actionCalledCount)
	assert.Equal(t, 1, errCalledCount)
}

func Test_DoWithRetries_ContextCanceledButWeIgnore(t *testing.T) {
	// given
	ctx, cancel := context.WithCancel(context.Background())
	actionCalledCount := 0
	action := func() (error, bool) {
		actionCalledCount++
		cancel()
		return fmt.Errorf("something went wrong"), true
	}

	errCalledCount := 0
	onError := func(err error) {
		errCalledCount++
		fmt.Println("Got an error: ", err.Error())
	}

	tr := NewRetrier(
		WithOnError(onError),
		WithIgnoreCtx(true),
	)

	// when
	err := tr.Do(ctx, action)

	// then
	require.Error(t, err)
	assert.EqualError(t, ErrNotSuccessful, err.Error())
	assert.Equal(t, 3, actionCalledCount)
	assert.Equal(t, 3, errCalledCount)
}