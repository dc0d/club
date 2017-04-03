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

	"github.com/dc0d/goroutines"
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
	werr := goroutines.New().
		WaitStart().
		WaitGo(timeout).
		Go(func() {
			AppPool.Wait()
		})
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
