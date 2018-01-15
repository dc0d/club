package om

import (
	"fmt"
	"testing"

	// "github.com/dc0d/club/om"
	"github.com/stretchr/testify/assert"
)

func Test01(t *testing.T) {
	assert := assert.New(t)
	o := New(nil)
	for i := 1; i <= 100; i++ {
		is := fmt.Sprintf("%03d", i)
		o.Put(is, i)
		assert.Equal(len(o._map), len(o._order))
	}
	for i := 1; i <= 100; i++ {
		is := fmt.Sprintf("%03d", i)
		v, ok := o.Get(is)
		assert.True(ok)
		assert.Equal(i, v)
	}
	itr := o.ItrFn()
	c := 1
	for k, v, ok := itr(); ok; k, v, ok = itr() {
		is := fmt.Sprintf("%03d", c)
		assert.Equal(is, k)
		assert.Equal(c, v)
		c++
	}
	for i := 1; i <= 100; i++ {
		is := fmt.Sprintf("%03d", i)
		o.Del(is)
		v, ok := o.Get(is)
		assert.False(ok)
		assert.Nil(v)
	}
	assert.Equal(0, len(o._map))
	assert.Equal(0, len(o._order))
}
