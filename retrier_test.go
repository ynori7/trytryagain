package trytryagain

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DoWithRetries_SuccessAfterSomeErrors(t *testing.T) {
	// given
	actionCalledCount := 0
	action := func() (error, bool) {
		actionCalledCount++
		if actionCalledCount > 2 {
			return nil, false
		}
		return fmt.Errorf("something went wrong"), true
	}

	errCalledCount := 0
	onError := func(err error) {
		errCalledCount++
		fmt.Println("Got an error: ", err.Error())
	}

	tr := NewRetrier(
		WithOnError(onError),
		WithBackoff(exponentialBackoff),
	)

	// when
	err := tr.Do(context.Background(), action)

	// then
	assert.NoError(t, err)
	assert.Equal(t, 3, actionCalledCount)
	assert.Equal(t, 2, errCalledCount)
}

func Test_DoWithRetries_ImmediateSuccess(t *testing.T) {
	// given
	actionCalledCount := 0
	action := func() (error, bool) {
		actionCalledCount++
		return nil, false
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
	err := tr.Do(context.Background(), action)

	// then
	assert.NoError(t, err)
	assert.Equal(t, 1, actionCalledCount)
	assert.Equal(t, 0, errCalledCount)
}

func Test_DoWithRetries_Failure(t *testing.T) {
	// given
	actionCalledCount := 0
	action := func() (error, bool) {
		actionCalledCount++
		return fmt.Errorf("something went wrong"), true
	}

	errCalledCount := 0
	onError := func(err error) {
		errCalledCount++
		fmt.Println("Got an error: ", err.Error())
	}

	tr := NewRetrier(
		WithMaxAttempts(4),
		WithOnError(onError),
	)

	// when
	err := tr.Do(context.Background(), action)

	// then
	require.Error(t, err)
	assert.EqualError(t, ErrNotSuccessful, err.Error())
	assert.Equal(t, 4, actionCalledCount)
	assert.Equal(t, 4, errCalledCount)
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
