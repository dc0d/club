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

// LoggerFunc wraps a func(args ...interface{}) in a Logger interface
type LoggerFunc func(args ...interface{})

func (lf LoggerFunc) Error(args ...interface{})                 { lf(args...) }
func (lf LoggerFunc) Info(args ...interface{})                  { lf(args...) }
func (lf LoggerFunc) Errorf(format string, args ...interface{}) { lf(fmt.Sprintf(format, args...)) }
func (lf LoggerFunc) Infof(format string, args ...interface{})  { lf(fmt.Sprintf(format, args...)) }

func init() {
	if log == nil {
		log = LoggerFunc(internalog.Print)
	}
}

// SetLogger replaces logger for this package
func SetLogger(logger Logger) { log = logger }

// GetLogger .
func GetLogger() Logger { return log }
