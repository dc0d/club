// Package logger provides a buffalo-compatible wrapper for zap
package logger

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/dc0d/club/group"
	"github.com/dc0d/config"
	"github.com/dc0d/config/hclconfig"
	"github.com/gobuffalo/buffalo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	confFile = "logger.conf"
)

// Conf .
type Conf struct {
	LogDir string
	Debug  bool
}

// Logger .
type Logger struct {
	*zap.SugaredLogger
}

// Printf .
func (l *Logger) Printf(t string, rest ...interface{}) {
	l.Infof(t, rest...)
}

// WithField .
func (l *Logger) WithField(k string, v interface{}) buffalo.Logger {
	return &Logger{l.With(k, v)}
}

// WithFields .
func (l *Logger) WithFields(kv map[string]interface{}) buffalo.Logger {
	var pairs []interface{}
	for k, v := range kv {
		pairs = append(pairs, k, v)
	}
	return &Logger{l.With(pairs...)}
}

// New .
func New(conf ...Conf) *Logger {
	var cnf Conf
	if len(conf) > 0 {
		cnf = conf[0]
	} else {
		fp, err := config.RelativeResource(confFile)
		if err != nil {
			// TODO:
		} else {
			if err := hclconfig.New().Load(&cnf, fp); err != nil {
				// TODO:
			}
		}
	}
	ld := os.TempDir()
	if cnf.LogDir != "" {
		ld = cnf.LogDir
	}

	ljdbg := ioutil.Discard
	if cnf.Debug {
		ljdbg = &lumberjack.Logger{
			Filename:   filepath.Join(ld, filepath.Base(os.Args[0])+"-debug.log"),
			MaxSize:    1, // megabytes
			MaxBackups: 3,
			MaxAge:     18, //days
		}
	}

	ljinfo := &lumberjack.Logger{
		Filename:   filepath.Join(ld, filepath.Base(os.Args[0])+"-info.log"),
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     18, //days
	}

	ljerr := &lumberjack.Logger{
		Filename:   filepath.Join(ld, filepath.Base(os.Args[0])+"-err.log"),
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     18, //days
	}

	priorityDbg := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.DebugLevel
	})
	priorityInfo := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel && lvl > zapcore.DebugLevel
	})
	priorityErr := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	topicDbg := zapcore.AddSync(ljdbg)
	topicInfo := zapcore.AddSync(ljinfo)
	topicErr := zapcore.AddSync(ljerr)

	consoleDbg := zapcore.Lock(os.Stdout)
	consoleInfo := zapcore.Lock(os.Stdout)
	consoleErr := zapcore.Lock(os.Stderr)

	var enccfg zapcore.EncoderConfig
	{
		pencconf := zap.NewProductionEncoderConfig()
		pencconf.EncodeTime = func(t time.Time, es zapcore.PrimitiveArrayEncoder) {
			es.AppendString(t.Format(time.RFC3339))
		}
		enccfg = pencconf
	}
	defaultEncoder := zapcore.NewJSONEncoder(enccfg)

	core := zapcore.NewTee(
		zapcore.NewCore(defaultEncoder, topicDbg, priorityDbg),
		zapcore.NewCore(defaultEncoder, consoleDbg, priorityDbg),
		zapcore.NewCore(defaultEncoder, topicInfo, priorityInfo),
		zapcore.NewCore(defaultEncoder, consoleInfo, priorityInfo),
		zapcore.NewCore(defaultEncoder, topicErr, priorityErr),
		zapcore.NewCore(defaultEncoder, consoleErr, priorityErr))

	_logger := zap.New(core, zap.AddCaller())
	loggerInstance := _logger.Sugar()
	zap.RedirectStdLog(_logger)

	cl := group.Child()
	cl.Add(1)
	go func() {
		defer cl.Done()
		<-cl.Stopped()
		loggerInstance.Sync()
		ljinfo.Close()
		ljerr.Close()
	}()

	return &Logger{loggerInstance}
}
