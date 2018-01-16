package om

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmoke(t *testing.T) {
	assert := assert.New(t)
	o := New(nil)
	for i := 1; i <= 100; i++ {
		is := fmt.Sprintf("%03d", i)
		o.Put(is, i)
		assert.Equal(len(o._map), len(o._order))
		assert.Equal(i, len(o._order))
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

func TestPutOrder(t *testing.T) {
	assert := assert.New(t)
	o := New(nil)
	for i := 1; i <= 100; i++ {
		for j := 1; j <= 100; j++ {
			is := fmt.Sprintf("%03d", i)
			o.Put(is, i)
			assert.Equal(len(o._map), len(o._order))
			assert.Equal(i, len(o._order))
			var buffer index = make([]keyT, len(o._order))
			copy(buffer, o._order)
			sort.Sort(buffer)
			for k := 0; k < len(buffer); k++ {
				assert.Equal(buffer[k], o._order[k])
			}
		}
	}
}

func TestDelOrder(t *testing.T) {
	assert := assert.New(t)
	o := New(nil)
	for i := 1; i <= 100; i++ {
		is := fmt.Sprintf("%03d", i)
		o.Put(is, i)
	}

	for i := 1; i <= 100; i++ {
		for j := 1; j <= 100; j++ {
			is := fmt.Sprintf("%03d", i)
			o.Del(is)
			assert.Equal(len(o._map), len(o._order))
			assert.Equal(100-i, len(o._order))
			var buffer index = make([]keyT, len(o._order))
			copy(buffer, o._order)
			sort.Sort(buffer)
			for k := 0; k < len(buffer); k++ {
				assert.Equal(buffer[k], o._order[k])
			}
		}
	}
}

func TestDelOrder2(t *testing.T) {
	assert := assert.New(t)
	o := New(nil)
	for i := 1; i <= 100; i++ {
		is := fmt.Sprintf("%03d", i)
		o.Put(is, i)
	}

	for i := 1; i <= 100; i++ {
		for j := 1; j <= 100; j++ {
			is := fmt.Sprintf("%03d", i)
			o.Del(is)
			assert.Equal(len(o._map), len(o._order))
			assert.Equal(99, len(o._order))
			var buffer index = make([]keyT, len(o._order))
			copy(buffer, o._order)
			sort.Sort(buffer)
			for k := 0; k < len(buffer); k++ {
				assert.Equal(buffer[k], o._order[k])
			}
			o.Put(is, i)
		}
	}
}

func TestItr(t *testing.T) {
	assert := assert.New(t)
	o := New(nil)
	for i := 1; i <= 1000; i++ {
		is := fmt.Sprintf("%04d", i)
		o.Put(is, i)
		assert.Equal(len(o._map), len(o._order))
		assert.Equal(i, len(o._order))
	}
	itr := o.ItrFn()
	c := 1
	var (
		prev1, prev2 keyT
	)
	for k, v, ok := itr(); ok; k, v, ok = itr() {
		is := fmt.Sprintf("%04d", c)
		assert.Equal(is, k)
		assert.Equal(c, v)
		k := k
		assert.Condition(func() bool { return prev2 <= prev1 && prev1 < k })
		c++
		prev1, prev2 = k, prev1
	}
}

func TestItrDesc(t *testing.T) {
	assert := assert.New(t)
	o := New(nil)
	for i := 1; i <= 1000; i++ {
		is := fmt.Sprintf("%04d", i)
		o.Put(is, i)
		assert.Equal(len(o._map), len(o._order))
		assert.Equal(i, len(o._order))
	}
	itr := o.ItrFn(true)
	c := 1000
	var (
		prev1, prev2 keyT
	)
	for k, v, ok := itr(); ok; k, v, ok = itr() {
		is := fmt.Sprintf("%04d", c)
		assert.Equal(is, k)
		assert.Equal(c, v)
		k := k
		assert.Condition(func() bool {
			if prev2 == "" {
				return true
			}
			if prev1 == "" {
				return true
			}
			return prev2 > prev1 && prev1 > k
		})
		c--
		prev1, prev2 = k, prev1
	}
	assert.Equal(len(o._map), len(o._order))
	assert.Equal(1000, len(o._order))
}

func BenchmarkPut(b *testing.B) {
	o := New(nil)
	for n := 0; n < b.N; n++ {
		o.Put("k", 1)
	}
}

func BenchmarkPutGet(b *testing.B) {
	o := New(nil)
	for n := 0; n < b.N; n++ {
		o.Put("k", 1)
		o.Get("k")
	}
}

func BenchmarkDel(b *testing.B) {
	o := New(nil)
	for i := 1; i <= 1000; i++ {
		is := fmt.Sprintf("%04d", i)
		o.Put(is, i)
	}
	for n := 0; n < b.N; n++ {
		// b.StopTimer()
		// o.Put("095", 95)
		// b.StartTimer()
		o.Del("0995")
	}
}

// check implementation
var _ sort.Interface = index{}
