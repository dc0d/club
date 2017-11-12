// Package singletonzaplogger is most likely a bad practice.
package singletonzaplogger

import (
	"os"
	"path/filepath"
	"time"

	singletongroup "github.com/dc0d/club/singleton-group"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

//-----------------------------------------------------------------------------

// Singleton returns the singleton instance of logger
func Singleton() Logger { return loggerInstance }

// Logger interface of *zap.SugaredLogger
type Logger interface {
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})

	Named(name string) *zap.SugaredLogger
	With(args ...interface{}) *zap.SugaredLogger
}

//-----------------------------------------------------------------------------

var (
	loggerInstance *zap.SugaredLogger

	grp = singletongroup.Singleton()

	ljinfo = &lumberjack.Logger{
		Filename:   "/tmp/" + filepath.Base(os.Args[0]) + "-info.log",
		MaxSize:    3, // megabytes
		MaxBackups: 6,
		MaxAge:     18, //days
	}

	ljerr = &lumberjack.Logger{
		Filename:   "/tmp/" + filepath.Base(os.Args[0]) + "-err.log",
		MaxSize:    3, // megabytes
		MaxBackups: 6,
		MaxAge:     18, //days
	}
)

//-----------------------------------------------------------------------------

func init() {
	priorityErr := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	priorityInfo := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	topicErr := zapcore.AddSync(ljerr)
	topicInfo := zapcore.AddSync(ljinfo)

	consoleErr := zapcore.Lock(os.Stderr)
	consoleInfo := zapcore.Lock(os.Stdout)

	var enccfg zapcore.EncoderConfig
	{
		pencconf := zap.NewProductionEncoderConfig()
		pencconf.EncodeTime = func(t time.Time, es zapcore.PrimitiveArrayEncoder) {
			es.AppendString(t.Format(time.RFC3339))
		}
		// pencconf.EncodeLevel = func(lvl zapcore.Level, arrenc zapcore.PrimitiveArrayEncoder) {
		// 	switch {
		// 	case lvl >= zapcore.ErrorLevel:
		// 		arrenc.AppendString(color.New(color.FgHiRed).Sprintf("%v", lvl))
		// 	case lvl == zapcore.WarnLevel:
		// 		arrenc.AppendString(color.New(color.FgHiYellow).Sprintf("%v", lvl))
		// 	default:
		// 		arrenc.AppendString(color.New(color.FgHiBlue).Sprintf("%v", lvl))
		// 	}
		// }
		enccfg = pencconf
	}
	defaultEncoder := zapcore.NewJSONEncoder(enccfg)

	core := zapcore.NewTee(
		zapcore.NewCore(defaultEncoder, topicErr, priorityErr),
		zapcore.NewCore(defaultEncoder, consoleErr, priorityErr),
		zapcore.NewCore(defaultEncoder, topicInfo, priorityInfo),
		zapcore.NewCore(defaultEncoder, consoleInfo, priorityInfo))

	_logger := zap.New(core, zap.AddCaller())
	loggerInstance = _logger.Sugar()
	zap.RedirectStdLog(_logger)
	go func() {
		ctx, wg := grp.Group()
		wg.Add(1)
		defer wg.Done()
		<-ctx.Done()
		loggerInstance.Sync()
		ljinfo.Close()
		ljerr.Close()
	}()

	<-time.After(time.Millisecond * 100)
}

//-----------------------------------------------------------------------------
