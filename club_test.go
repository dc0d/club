package club

import (
	"testing"
)

func Test01(t *testing.T) {
	defer TimerScope("")()
	// // zap.NewProduction(zap.)
	// logger, _ := zap.NewProduction()
	// logger.Sugar()
	// defer logger.Sync()
	// logger.Info("failed to fetch URL",
	// 	// Structured context as strongly typed Field values.
	// 	zap.Int("attempt", 3),
	// 	zap.Duration("backoff", time.Second),
	// )
}
