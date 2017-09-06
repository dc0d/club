package errors

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

//-----------------------------------------------------------------------------

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
	buff.Reset()
	bufferPool.Put(buff)
}

//-----------------------------------------------------------------------------

// ErrorCollection is an error type representing a collection of errors
type ErrorCollection []error

func (x ErrorCollection) Error() string {
	if len(x) == 0 {
		return ""
	}

	buff := getBuffer()
	defer putBuffer(buff)

	for _, ve := range x {
		if ve == nil {
			continue
		}
		buff.WriteString(" [" + ve.Error() + "]")
	}
	res := strings.TrimSpace(buff.String())

	return res
}

//-----------------------------------------------------------------------------

// ErrString a string that inplements the error interface
type ErrString string

func (v ErrString) Error() string { return string(v) }

// Errorf value type (string) error
func Errorf(format string, a ...interface{}) error {
	return ErrString(fmt.Sprintf(format, a...))
}

//-----------------------------------------------------------------------------

// Here(...) errors
var (
	ErrNotAvailable = Errorf("N/A")
)

// Here .
func Here(skip ...int) (funcName, fileName string, fileLine int, callerErr error) {
	sk := 1
	if len(skip) > 0 && skip[0] > 1 {
		sk = skip[0]
	}
	var pc uintptr
	var ok bool
	pc, fileName, fileLine, ok = runtime.Caller(sk)
	if !ok {
		callerErr = ErrNotAvailable
		return
	}
	fn := runtime.FuncForPC(pc)
	name := fn.Name()
	ix := strings.LastIndex(name, ".")
	if ix > 0 && (ix+1) < len(name) {
		name = name[ix+1:]
	}
	funcName = name
	nd, nf := filepath.Split(fileName)
	fileName = filepath.Join(filepath.Base(nd), nf)
	return
}

//-----------------------------------------------------------------------------

// ErrorCallerf creates a string error which containes the info about location of error
func ErrorCallerf(format string, a ...interface{}) error {
	var name string
	funcName, fileName, fileLine, err := Here(2)
	if err != nil {
		name = "N/A"
	} else {
		name = fmt.Sprintf("%s:%02d %s()", fileName, fileLine, funcName)
	}
	return Errorf(name+": "+format, a...)
}

//-----------------------------------------------------------------------------
