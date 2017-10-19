package errgroup

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dc0d/club/errgo/sentinel"
	"github.com/stretchr/testify/assert"
)

func TestSampleClosure(t *testing.T) {
	grp, gctx := WithContext(context.Background())

	actions := &actionCollection{Context: gctx}

	for i := 0; i < 100; i++ {
		i := i
		grp.Go(actions.action(i))
	}

	assert.Equal(t, ErrExpectedSample, grp.Wait())
}

type actionCollection struct {
	// shared closure for goroutines
	context.Context
	count int64
}

func (ac *actionCollection) action(i int) func() error {
	return func() error {
		for {
			select {
			case <-ac.Done():
				return nil
			case <-time.After(time.Millisecond * 5):
				atomic.AddInt64(&ac.count, int64(1))
			}
			last := atomic.LoadInt64(&ac.count)
			if last > 10 {
				panic(ErrExpectedSample)
			}
		}
	}
}

var ErrExpectedSample = sentinel.Errorf("SHOULE_STOP")
