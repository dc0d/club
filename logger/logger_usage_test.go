package logger_test

import (
	"os"
	"testing"

	"github.com/dc0d/club/group"
	"github.com/dc0d/club/logger"
)

var (
	_logger *logger.Logger
)

func init() {
	_logger = logger.New(logger.Conf{LogDir: os.TempDir(), Debug: true})
}

func Test01(t *testing.T) {
	_logger.Debug("DEBUG")
	_logger.Info("Info")
	_logger.Warn("WARN")
	_logger.Error("ERROR")
	group.Supervisor().Stop()
	group.Supervisor().Wait()
}
