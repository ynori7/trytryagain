package trytryagain

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IsCanceledContextError(t *testing.T) {
	// given
	testcases := map[string]bool{
		`{"source":"","category":408,"code":"context_canceled","message":"Context canceled","data":"","retriable":true}`: true,
		`Get "https://www.blah.com/something": context canceled`:                                                         true,
		`context deadline exceeded`: true,
		`not found`:                 false,
	}

	for testcase, expected := range testcases {
		// when
		actual := IsCanceledContextError(fmt.Errorf(testcase))

		// then
		assert.Equal(t, expected, actual)
	}
}

func Test_IsContextDone(t *testing.T) {
	// given
	ctx, cancel := context.WithCancel(context.Background())

	// when
	actual := IsContextDone(ctx)

	// then
	assert.False(t, actual)

	// when
	cancel()
	actual = IsContextDone(ctx)

	//  then
	assert.True(t, actual)
}
