package clubaux

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dc0d/errgo/sentinel"
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
	errNotAvailable = sentinel.Errorf("N/A")
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

// TruncateHour .
func TruncateHour(t time.Time) (res time.Time) {
	res = t.Truncate(time.Minute * 30)
	if res.Minute() > 0 {
		res = res.Add(-1 * time.Minute).Truncate(time.Minute * 30)
	}
	return
}

//-----------------------------------------------------------------------------
