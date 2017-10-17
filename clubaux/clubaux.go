package clubaux

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dc0d/club/errors"
)

//-----------------------------------------------------------------------------

// OnSignal runs function on receiving the OS signal
func OnSignal(f func(), sig ...os.Signal) {
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
			syscall.SIGTSTP,
			syscall.SIGKILL)
	}
	go func() {
		<-sigc
		f()
	}()
}

var (
	errNotAvailable = errors.Errorf("N/A")
)

//-----------------------------------------------------------------------------

// Chain .
func Chain(steps ...func() error) (chainerr error) {
	for _, v := range steps {
		chainerr = v()
		if chainerr != nil {
			return
		}
	}
	return
}

//-----------------------------------------------------------------------------
