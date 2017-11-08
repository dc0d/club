package wg_test

import (
	"testing"
	"time"

	"github.com/dc0d/club/wg"
	"github.com/stretchr/testify/assert"
)

var singleWG = wg.New()

func Test01(t *testing.T) {
	assert := assert.New(t)
	singleWG.Add(1)
	go func() {
		defer singleWG.Done()
	}()
	<-time.After(time.Microsecond)
	err := singleWG.Wait()
	assert.NoError(err)
}

func Test02(t *testing.T) {
	assert := assert.New(t)
	singleWG.Add(1)
	go func() {
		defer singleWG.Done()
		time.Sleep(time.Second)
	}()
	err := singleWG.Wait(time.Millisecond)
	assert.Error(err)
	assert.Equal(err, wg.ErrTimeout)
}
