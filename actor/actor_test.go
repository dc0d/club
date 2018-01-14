package actor_test

import (
	"testing"
	"time"

	"github.com/dc0d/club/actor"
	"github.com/stretchr/testify/assert"
)

var (
	mailboxSize     = actor.MailboxSize
	numberOfRetries = actor.NumberOfRetries
	onError         = actor.OnError
	period          = actor.Period
)

func Test01(t *testing.T) {
	assert := assert.New(t)

	errors := make(chan error, 10)
	actor.Start(func(mailbox <-chan interface{}) error {
		panic(10)
	},
		mailboxSize(10),
		period(time.Millisecond*10),
		numberOfRetries(3),
		onError(func(e error) { errors <- e }))

	var c int
FOR01:
	for i := 1; i <= 10; i++ {
		select {
		case <-errors:
		case <-time.After(time.Millisecond * 100):
			break FOR01
		}
		c++
	}
	assert.Equal(3, c)
}
