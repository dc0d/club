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

// Supervise runs in sync, use as "go Supervise(...)",
// takes care of restarts (in case of panic or error)
func Supervise(action func() error, intensity int, period ...time.Duration) {
	dt := time.Second * 3
	if len(period) > 0 && period[0] > 0 {
		dt = period[0]
	}
	for intensity != 0 {
		if intensity > 0 {
			intensity--
		}
		if err := Run(action); err != nil {
			log.Errorf("supervised %v", err)
			if intensity != 0 {
				time.Sleep(dt)
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
