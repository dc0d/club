package timerscope

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/dc0d/club/errors"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func getBuffer() *bytes.Buffer {
	buff := bufferPool.Get().(*bytes.Buffer)
	return buff
}

func putBuffer(buff *bytes.Buffer) {
	buff.Reset() // ooch!
	bufferPool.Put(buff)
}

type option struct {
	name    string
	opCount int
}

// Option .
type Option func(option) option

// Name .
func Name(name string) Option {
	return func(opt option) option {
		opt.name = name
		return opt
	}
}

// OpCount .
func OpCount(opCount int) Option {
	return func(opt option) option {
		opt.opCount = opCount
		return opt
	}
}

// TimerScope .
func TimerScope(options ...Option) (name string, onExit func() string) {
	var opt option
	for _, v := range options {
		opt = v(opt)
	}
	if opt.name == "" {
		funcName, fileName, fileLine, err := errors.Here(2)
		if err != nil {
			opt.name = "N/A"
		} else {
			opt.name = fmt.Sprintf("%s:%02d %s()", fileName, fileLine, funcName)
		}
	}
	name = opt.name
	start := time.Now()
	onExit = func() (logExit string) {
		buf := getBuffer()
		defer putBuffer(buf)
		defer func() {
			logExit = string(buf.Bytes())
		}()

		elapsed := time.Now().Sub(start)
		fmt.Fprintf(buf, "%s took %v ", name, elapsed)

		N := opt.opCount
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

		return
	}

	return
}
