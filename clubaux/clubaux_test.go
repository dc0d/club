package clubaux

import (
	"errors"
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
	err := LoadHCL(&sampleConf)
	assert.Nil(t, err)
	assert.Equal(t, "OK", sampleConf.Info)
	assert.Equal(t, "OK", sampleConf.Sample.SubCommand.Param)
}

var _ error = ErrorCollection(nil)

func Test02(t *testing.T) {
	var x ErrorCollection
	x = append(x, errors.New("ERR 1"))
	x = append(x, errors.New("ERR 2"))

	var err error = x
	assert.Equal(t, "[ERR 1] [ERR 2]", err.Error())
}
