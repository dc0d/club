package hclconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleConf struct {
	Info string

	Sample struct {
		SubCommand struct {
			Param string
		}
	}
}

func Test01(t *testing.T) {
	err := New().Load(&sampleConf)
	assert.Nil(t, err)
	assert.Equal(t, "OK", sampleConf.Info)
	assert.Equal(t, "OK", sampleConf.Sample.SubCommand.Param)
}
