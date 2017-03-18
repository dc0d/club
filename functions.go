package club

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

//-----------------------------------------------------------------------------

func init() {
	AppCtx, AppCancel = context.WithCancel(context.Background())
	CallOnSignal(func() { AppCancel() })
	AppPool = NewGroup()
}

//-----------------------------------------------------------------------------

// Finit use to prevent app from stopping/exiting
func Finit(timeout time.Duration, nap ...time.Duration) {
	<-AppCtx.Done()
	werr := WaitGo(func() {
		AppPool.Wait()
	}, timeout)
	if werr != nil {
		log.Println(`error:`, werr)
	}
	s := time.Millisecond * 300
	if len(nap) > 0 {
		s = nap[0]
	}
	if s <= 0 {
		s = time.Millisecond * 300
	}
	<-time.After(s)
}

//-----------------------------------------------------------------------------

// WaitGo waits for f to complete (in a goroutine)
// or times out (returning ErrTimeout)
func WaitGo(f func(), timeout ...time.Duration) error {
	funcDone := make(chan struct{})
	go func() {
		defer close(funcDone)
		f()
	}()

	var delay time.Duration
	if len(timeout) > 0 {
		delay = timeout[0]
	}

	if delay <= 0 {
		<-funcDone

		return nil
	}

	select {
	case <-time.After(delay):
		return ErrTimeout
	case <-funcDone:
	}

	return nil
}

//-----------------------------------------------------------------------------

// WaitStart starts a goroutine and wait for it to start, and after goroutine
// started, it returns.
func WaitStart(f func()) {
	started := make(chan struct{})
	go func() {
		close(started)
		f()
	}()
	<-started
}

//-----------------------------------------------------------------------------

// Recover recovers from panic and returns error,
// or returns the provided error
func Recover(f func() error, verbose ...bool) (err error) {
	vrb := false
	if len(verbose) > 0 {
		vrb = verbose[0]
	}

	defer func() {
		if e := recover(); e != nil {
			if !vrb {
				err = Error(fmt.Sprintf("%v", e))
				return
			}

			trace := make([]byte, 1<<16)
			n := runtime.Stack(trace, true)

			s := fmt.Sprintf("error: %v, n: %d, traced: %s", e, n, trace[:n])
			err = Error(s)
		}
	}()

	err = f()

	return
}

//-----------------------------------------------------------------------------

// CallOnSignal calls function on os signal
func CallOnSignal(f func(), sig ...os.Signal) {
	if f == nil {
		return
	}
	sigc := make(chan os.Signal, 1)
	if len(sig) > 0 {
		signal.Notify(sigc, sig...)
	} else {
		signal.Notify(sigc,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
			syscall.SIGSTOP,
			syscall.SIGABRT,
			syscall.SIGTSTP,
			syscall.SIGKILL)
	}
	go func() {
		<-sigc
		f()
	}()
}

//-----------------------------------------------------------------------------
