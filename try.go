package try

import "errors"

var (
	MaxRetries           = 10                                 // MaxRetries is the maximum number of retries before bailing.
	ErrMaxRetriesReached = errors.New("exceeded retry limit") // ErrMaxRetriesReached is returned when the function has failed to return true before MaxRetries is reached.
)

// Func is the function to be retried.
type Func func(attempt int) (retry bool, err error)

// Do keeps trying the function until returns false, or no error is returned.
func Do(fn Func) (err error) {
	var cont bool

	attempt := 1

	for {
		cont, err = fn(attempt)
		if !cont || err == nil {
			break
		}

		attempt++

		if attempt > MaxRetries {
			return ErrMaxRetriesReached
		}
	}

	return err
}
