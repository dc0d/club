package club

import (
	"fmt"
	"time"
)

//-----------------------------------------------------------------------------

// ErrString a string that inplements the error interface
type ErrString string

func (v ErrString) Error() string { return string(v) }

// Errorf value type (string) error
func Errorf(format string, a ...interface{}) error {
	return ErrString(fmt.Sprintf(format, a...))
}

//-----------------------------------------------------------------------------

// Supervise a helper method, runs in sync, use as "go Supervise(...)",
// takes care of restarts (in case of panic or error),
// can be used as a SimpleOneForOne supervisor, an intensity < 0 means restart forever
func Supervise(
	action func() error,
	intensity int,
	period time.Duration,
	onError ...func(error)) {
	if intensity != 1 && period <= 0 {
		period = time.Second * 5
	}

	for intensity != 0 {
		if intensity > 0 {
			intensity--
		}
		if err := Run(action); err != nil {
			if len(onError) > 0 && onError[0] != nil {
				onError[0](err)
			}
			if intensity != 0 {
				time.Sleep(period)
			}
		} else {
			break
		}
	}
}

// Run calls the function, does captures panics
func Run(action func() error) (errrun error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				errrun = err
				return
			}
			errrun = Errorf("UNKNOWN: %v", e)
		}
	}()
	return action()
}

//-----------------------------------------------------------------------------
