package clubaux

import (
	"context"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTruncateHour(t *testing.T) {
	t.SkipNow()
	assert := assert.New(t)

	src := rand.NewSource(time.Now().Unix())
	r := rand.New(src)

	for i := 0; i < 100; i++ {
		step := time.Second * time.Duration(r.Intn(3600)+1)

		end := time.Now().Add(24 * 30 * 12 * time.Hour)
		current := time.Now()
		for current.Before(end) {
			current = current.Add(step)

			th := TruncateHour(current)
			assert.Zero(th.Minute())
			assert.Equal(current.Hour(), th.Hour())
		}
	}
}

func TestThrottle(t *testing.T) {
	assert := assert.New(t)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Millisecond*100))
	defer cancel()

	count := 0
	ttt := Throttle(ctx, 30, time.Millisecond*10)
OUT1:
	for {
		select {
		case <-ctx.Done():
			break OUT1
		case _, ok := <-ttt:
			if ok {
				count++
			} else {
				break OUT1
			}
		}
	}
	assert.Condition(func() bool {
		diff := math.Abs(float64(300 - count))
		return diff <= 100
	})
}

func TestThrottleNoContext(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Millisecond*100))
	defer cancel()

	count := 0
	ttt := Throttle(nil, 30, time.Millisecond*10)
OUT1:
	for {
		select {
		case <-ctx.Done():
			break OUT1
		case _, ok := <-ttt:
			if ok {
				count++
			} else {
				break OUT1
			}
		}
	}
	assert.Condition(func() bool {
		diff := math.Abs(float64(300 - count))
		return diff <= 30
	})
}
