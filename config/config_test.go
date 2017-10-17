package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test01(t *testing.T) {
	assert := assert.New(t)

	p, err := RelativeResource("functions.go")
	assert.NoError(err)
	assert.Contains(p, "/src/github.com/dc0d/club/config/functions.go")
}
