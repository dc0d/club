package tgroup

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test01(t *testing.T) {
	assert := assert.New(t)

	octx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var c int64

	g := New(octx)
	_, wg := g.Group()

	N := 10
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			atomic.AddInt64(&c, 1)
		}()
	}
	go func() { cancel() }()
	err := WaitFinish(g, time.Second)
	assert.NoError(err)
	assert.Equal(int64(N), c)
}

func Test02(t *testing.T) {
	assert := assert.New(t)

	octx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var c int64

	g := New(octx)
	ctx, wg := g.Group()

	N := 10
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}
			atomic.AddInt64(&c, 1)
		}()
	}
	cancel()
	err := WaitFinish(g, time.Second)
	assert.NoError(err)
	assert.Condition(func() bool { return c < int64(N) })
}
