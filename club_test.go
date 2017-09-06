package club

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dc0d/club/errors"
	"github.com/stretchr/testify/assert"
)

func TestByteSize(t *testing.T) {
	var n ByteSize = 2048
	assert.Equal(t, "2.00 KB", fmt.Sprint(n))
}

func Test01(t *testing.T) {
	SetLogger(LoggerFunc(func(...interface{}) {}))
	var sum int64
	Supervise(func() error {
		atomic.AddInt64(&sum, 1)
		return errors.Errorf("DUMMY")
	}, 3, time.Millisecond*50)
	assert.Equal(t, int64(3), sum)
}

func Test02(t *testing.T) {
	SetLogger(LoggerFunc(func(...interface{}) {}))
	var sum int64
	Supervise(func() error {
		atomic.AddInt64(&sum, 1)
		panic(errors.Errorf("DUMMY"))
	}, 3, time.Millisecond*50)
	assert.Equal(t, int64(3), sum)
}

func Test03(t *testing.T) {
	SetLogger(LoggerFunc(func(...interface{}) {}))
	var sum int64
	Supervise(func() error {
		atomic.AddInt64(&sum, 1)
		return nil
	}, 3, time.Millisecond*50)
	assert.Equal(t, int64(1), sum)
}
