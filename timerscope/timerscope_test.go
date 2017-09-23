package timerscope

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test01(t *testing.T) {
	assert := assert.New(t)

	fst, sndFunc := TimerScope()
	snd := sndFunc()

	assert.Contains(fst, "timerscope/timerscope_test.go")
	assert.Contains(fst, "Test01()")
	assert.Contains(snd, "timerscope/timerscope_test.go")
	assert.Contains(snd, "Test01() took")
}

func Test02(t *testing.T) {
	assert := assert.New(t)

	fst, sndFunc := TimerScope(Name("sample"), OpCount(100))
	snd := sndFunc()

	assert.Equal(fst, "sample")
	assert.Contains(snd, "sample took")
	assert.Contains(snd, "op/sec")
	assert.Contains(snd, "nano-sec/op")
}

func Test03(t *testing.T) {
	t.SkipNow()
	fst, sndFunc := TimerScope(Name("sample"), OpCount(100))
	t.Log(fst)
	defer func() { t.Log(sndFunc()) }()

	time.Sleep(time.Second)
}
