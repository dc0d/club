package club

import (
	"fmt"
	internalog "log"
)

// Logger generic logger
type Logger interface {
	Error(args ...interface{})
	Info(args ...interface{})
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
}

var log Logger

type loggerFunc func(args ...interface{})

func (lf loggerFunc) Error(args ...interface{})                 { lf(args...) }
func (lf loggerFunc) Info(args ...interface{})                  { lf(args...) }
func (lf loggerFunc) Errorf(format string, args ...interface{}) { lf(fmt.Sprintf(format, args...)) }
func (lf loggerFunc) Infof(format string, args ...interface{})  { lf(fmt.Sprintf(format, args...)) }

func init() {
	if log == nil {
		log = loggerFunc(internalog.Print)
	}
}
