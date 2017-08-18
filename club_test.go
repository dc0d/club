package club

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteSize(t *testing.T) {
	var n ByteSize = 2048
	assert.Equal(t, "2.00 KB", fmt.Sprint(n))
}

func Test01(t *testing.T) {
	// defer TimerScope("")()
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
