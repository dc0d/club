package iniconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Common struct {
	WD  string `ini:"wd"`  // = .
	Tmp string `ini:"tmp"` // = /tmp
	Max int    `ini:"max"` // = 10
	Min int    `ini:"min"` // = 3
}

func Test01(t *testing.T) {
	assert := assert.New(t)

	c := new(Common)
	err := New().Load(c)
	assert.Nil(err)
	assert.Equal(".", c.WD)
	assert.Equal("/tmp", c.Tmp)
	assert.Equal(10, c.Max)
	assert.Equal(3, c.Min)
}
