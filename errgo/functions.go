package errgo

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dc0d/club/errgo/sentinel"
)

//-----------------------------------------------------------------------------

// ErrorMarkf .
func ErrorMarkf(format string, a ...interface{}) error {
	name := makeMarker(markerFormat)
	cause := sentinel.Errorf(format, a...)
	return &Marker{
		name: name,
		err:  cause,
	}
}

//-----------------------------------------------------------------------------

// MarkHere .
func MarkHere() (name string) { return makeMarker(markerFormat) }

//-----------------------------------------------------------------------------

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

func makeMarker(format string) (name string) {
	funcName, fileName, fileLine, err := Here(3)
	if err != nil {
		name = "N/A"
	} else {
		name = fmt.Sprintf(markerFormat, fileName, fileLine, funcName)
	}
	return name
}

//-----------------------------------------------------------------------------
