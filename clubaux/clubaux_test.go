package clubaux

import (
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
