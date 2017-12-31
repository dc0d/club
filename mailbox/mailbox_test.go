package mailbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Mailbox = &mailbox{}

func TestSmoke(t *testing.T) {
	assert := assert.New(t)
	mbox := New(&SliceStorage{})
	go func() { mbox.Send() <- 1 }()
	v := <-mbox.Receive()
	assert.Equal(1, v)
}

func TestCount(t *testing.T) {
	assert := assert.New(t)
	store := SliceStorage{}
	mbox := New(&store)
	N := 10
	go func() {
		for i := 0; i < N; i++ {
			mbox.Send() <- 1
		}
		mbox.Close()
	}()
	total := 0
OUT1:
	for {
		select {
		case v, ok := <-mbox.Receive():
			if !ok {
				break OUT1
			}
			assert.Equal(1, v)
			total++
		}
	}
	assert.Equal(N, total)
}

func TestItems(t *testing.T) {
	mbox := New(&SliceStorage{})
	N := 1000
	go func() {
		for i := 1; i <= N; i++ {
			i := i
			mbox.Send() <- i
		}
		mbox.Close()
	}()
	total := 0
OUT1:
	for {
		select {
		case v, ok := <-mbox.Receive():
			if !ok {
				break OUT1
			}
			total += v.(int)
		}
	}
	assert.Equal(t, 500500, total)
}
