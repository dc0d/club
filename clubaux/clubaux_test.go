package clubaux

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test01(t *testing.T) {
	rg := NewRegistry(nil, time.Millisecond*30)

	rg.Put(1, 1)
	v, ok := rg.Get(1)
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	rg.PutWithExpiration(2, 2, time.Millisecond*50)
	v, ok = rg.Get(2)
	assert.True(t, ok)
	assert.Equal(t, 2, v)
	<-time.After(time.Millisecond * 100)

	v, ok = rg.Get(2)
	assert.False(t, ok)
	assert.NotEqual(t, 2, v)
}

func Test02(t *testing.T) {
	rg := NewRegistry(nil, time.Millisecond*30)

	rg.Put(1, 1)
	v, ok := rg.Get(1)
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	rg.PutWithExpiration(1, 1, time.Millisecond*50, true)
	<-time.After(time.Millisecond * 40)
	v, ok = rg.Get(1)
	assert.True(t, ok)
	assert.Equal(t, 1, v)
	<-time.After(time.Millisecond * 10)
	v, ok = rg.Get(1)
	assert.True(t, ok)
	assert.Equal(t, 1, v)
	<-time.After(time.Millisecond * 10)
	v, ok = rg.Get(1)
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	<-time.After(time.Millisecond * 100)

	v, ok = rg.Get(1)
	assert.False(t, ok)
	assert.NotEqual(t, 1, v)
}

func BenchmarkGet01(b *testing.B) {
	rg := NewRegistry(nil, 0)
	for n := 0; n < b.N; n++ {
		rg.Get(1)
	}
}

func BenchmarkGet02(b *testing.B) {
	rg := NewRegistry(nil, 0)
	rg.Put(1, 1)
	for n := 0; n < b.N; n++ {
		rg.Get(1)
	}
}

func BenchmarkGet03(b *testing.B) {
	rg := NewRegistry(nil, 0)
	rg.PutWithExpiration(1, 1, time.Second)
	for n := 0; n < b.N; n++ {
		rg.Get(1)
	}
}

func BenchmarkPut01(b *testing.B) {
	rg := NewRegistry(nil, 0)
	for n := 0; n < b.N; n++ {
		rg.Put(1, 1)
	}
}

func BenchmarkCAS02(b *testing.B) {
	rg := NewRegistry(nil, 0)
	rg.Put(1, 1)
	for n := 0; n < b.N; n++ {
		rg.CAS(1, 2, func(interface{}) bool { return true })
	}
}
