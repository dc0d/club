package clubaux

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/dc0d/club"
	"github.com/dc0d/club/errors"
	"github.com/hashicorp/hcl"
)

//-----------------------------------------------------------------------------

var log = club.GetLogger()

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
			syscall.SIGABRT,
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

// TimerScope .
func TimerScope(name string, opCount ...int) func() {
	if name == "" {
		funcName, fileName, fileLine, err := errors.Here(2)
		if err != nil {
			name = "N/A"
		} else {
			name = fmt.Sprintf("%s:%02d %s()", fileName, fileLine, funcName)
		}
	}
	log.Info(name, " started")
	start := time.Now()
	return func() {
		buf := GetBuffer()
		defer PutBuffer(buf)
		defer func() {
			log.Info(string(buf.Bytes()) + "\n")
		}()

		elapsed := time.Now().Sub(start)
		fmt.Fprintf(buf, "%s took %v ", name, elapsed)
		if len(opCount) == 0 {
			return
		}

		N := opCount[0]
		if N <= 0 {
			return
		}

		E := float64(elapsed)
		FRC := E / float64(N)

		fmt.Fprintf(buf, "op/sec %.2f ", float64(N)/(E/float64(time.Second)))

		switch {
		case FRC > float64(time.Second):
			fmt.Fprintf(buf, "sec/op %.2f ", (E/float64(time.Second))/float64(N))
		case FRC > float64(time.Millisecond):
			fmt.Fprintf(buf, "milli-sec/op %.2f ", (E/float64(time.Millisecond))/float64(N))
		case FRC > float64(time.Microsecond):
			fmt.Fprintf(buf, "micro-sec/op %.2f ", (E/float64(time.Microsecond))/float64(N))
		default:
			fmt.Fprintf(buf, "nano-sec/op %.2f ", (E/float64(time.Nanosecond))/float64(N))
		}
	}
}

//-----------------------------------------------------------------------------

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

// GetBuffer .
func GetBuffer() *bytes.Buffer {
	buff := bufferPool.Get().(*bytes.Buffer)
	return buff
}

// PutBuffer .
func PutBuffer(buff *bytes.Buffer) {
	buff.Reset() // ooch!
	bufferPool.Put(buff)
}

//-----------------------------------------------------------------------------

// LoadHCL loads hcl conf file. default conf file names (if filePath not provided)
// in the same directory are <appname>.conf and if not fount app.conf
func LoadHCL(ptr interface{}, filePath ...string) error {
	var fp string
	if len(filePath) > 0 {
		fp = filePath[0]
	}
	if fp == "" {
		fp = _confFilePath()
	}
	cn, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	err = hcl.Unmarshal(cn, ptr)
	if err != nil {
		return err
	}

	return nil
}

func _confFilePath() string {
	appName := filepath.Base(os.Args[0])
	appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {

	}
	appConfName := fmt.Sprintf("%s.conf", appName)
	genericConfName := "app.conf"

	for _, vn := range []string{appConfName, genericConfName} {
		currentPath := filepath.Join(appDir, vn)
		if _, err := os.Stat(currentPath); err == nil {
			return currentPath
		}
	}

	for _, vn := range []string{appConfName, genericConfName} {
		wd, err := os.Getwd()
		if err != nil {
			continue
		}
		currentPath := filepath.Join(wd, vn)
		if _, err := os.Stat(currentPath); err == nil {
			return currentPath
		}
	}

	if _, err := os.Stat(appConfName); err == nil {
		return appConfName
	}

	return genericConfName
}

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
