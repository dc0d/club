package mailbox

import (
	"context"
	"testing"

	"github.com/dc0d/club/mailbox/rammailstorage"
	"github.com/stretchr/testify/assert"
)

func TestSmoke(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mbox := New(ctx, rammailstorage.New())
	go func() { mbox.Send() <- 1 }()
	v := <-mbox.Receive()
	assert.Equal(t, 1, v)
}

func TestCount(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mbox := New(ctx, rammailstorage.New())
	N := 1000
	go func() {
		for i := 0; i < N; i++ {
			mbox.Send() <- 1
		}
		defer cancel()
	}()
	total := 0
OUT1:
	for {
		select {
		case v, ok := <-mbox.Receive():
			if !ok {
				break OUT1
			}
			assert.Equal(t, 1, v)
			total++
		}
	}
	assert.Equal(t, N, total)
}

func TestItems(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mbox := New(ctx, rammailstorage.New())
	N := 1000
	go func() {
		for i := 1; i <= N; i++ {
			i := i
			mbox.Send() <- i
		}
		defer cancel()
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
