package errgo

import (
	"strings"
	"testing"

	"github.com/dc0d/club/errgo/sentinel"
	"github.com/stretchr/testify/assert"
)

func removeNumbers(s string) string {
	rep := strings.NewReplacer(
		"1", "_",
		"2", "_",
		"3", "_",
		"4", "_",
		"5", "_",
		"6", "_",
		"7", "_",
		"8", "_",
		"9", "_",
		"0", "_")
	return rep.Replace(s)
}

func TestMarker(t *testing.T) {
	assert := assert.New(t)

	m := NewMarker(sentinel.Errorf("BOO"))
	assert.Error(m)
	assert.Contains(removeNumbers(m.Error()), removeNumbers(MarkHere()))
}

func TestErrorMarkf(t *testing.T) {
	assert := assert.New(t)
	cause := sentinel.Errorf("BOO %d", 10)
	em := ErrorMarkf("BOO %d", 10)
	if c, ok := em.(interface {
		Cause() error
	}); ok {
		assert.Equal(cause, c.Cause())
	}
	assert.Contains(removeNumbers(em.Error()), removeNumbers(MarkHere()))
}
