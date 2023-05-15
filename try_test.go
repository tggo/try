package try_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tggo/try"
)

func TestTryExample(t *testing.T) {
	try.MaxRetries = 20

	SomeFunction := func() (string, error) {
		return "OK", nil
	}

	var value string
	err := try.Do(func(attempt int) (condition bool, err error) {
		value, err = SomeFunction()
		return attempt < 5, err // try 5 times
	})

	require.NoError(t, err, "error should be nil")
	require.Equal(t, "OK", value, "value should be OK")
}

func TestTryExampleWithSleep(t *testing.T) {
	try.MaxRetries = 20

	SomeFunction := func() (string, error) {
		return "OK", nil
	}

	var value string
	err := try.Do(func(attempt int) (condition bool, err error) {
		value, err = SomeFunction()
		if err != nil {
			time.Sleep(1 * time.Millisecond)
		}
		return attempt < 5, err // try 5 times
	})

	require.NoError(t, err, "error should be nil")
	require.Equal(t, "OK", value, "value should be OK")
}

func TestTryExamplePanic(t *testing.T) {
	SomeFunction := func() (string, error) {
		panic("something went badly wrong")
	}

	var value string

	err := try.Do(func(attempt int) (retry bool, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("panic: %v", r))
			}
		}()

		retry = attempt < 5 // try 5 times
		value, err = SomeFunction()

		return
	})

	require.Equal(t, "", value, "value should be empty")
	require.Equal(t, "panic: something went badly wrong", err.Error(), "error should be 'panic: something went badly wrong'")
}

func TestTryDoSuccessful(t *testing.T) {
	callCount := 0
	err := try.Do(func(attempt int) (bool, error) {
		callCount++
		return attempt < 5, nil
	})

	require.NoError(t, err, "error should be nil")
	require.Equal(t, 1, callCount, "callCount should be 1")
}

func TestTryDoFailed(t *testing.T) {
	theErr := errors.New("something went wrong")
	callCount := 0

	err := try.Do(func(attempt int) (bool, error) {
		callCount++
		return attempt < 5, theErr
	})

	require.Equal(t, theErr, err, "error should be theErr")
	require.Equal(t, 5, callCount, "callCount should be 5")
}

func TestTryPanics(t *testing.T) {
	theErr := errors.New("something went wrong")
	callCount := 0
	err := try.Do(func(attempt int) (retry bool, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("panic: %v", r))
			}
		}()

		retry = attempt < 5
		callCount++

		if attempt > 2 {
			panic("is so much")
		}
		err = theErr

		return
	})

	require.Equal(t, "panic: is so much", err.Error(), "error should 'is so much'")
	require.Equal(t, 5, callCount, "callCount should be 5")
}

func TestRetryLimit(t *testing.T) {
	err := try.Do(func(attempt int) (bool, error) {
		return true, errors.New("nope")
	})

	require.Error(t, err, "error should not be nil")
	require.Equal(t, try.ErrMaxRetriesReached, err)
}
