package actor

import (
	"time"

	"github.com/dc0d/retry"
)

type payload = interface{}

type options struct {
	mailboxSize     int
	numberOfRetries int
	onError         func(error)
	period          time.Duration
}

// MailboxSize .
func MailboxSize(v int) func(o *options) { return func(o *options) { o.mailboxSize = v } }

// NumberOfRetries .
func NumberOfRetries(v int) func(o *options) { return func(o *options) { o.numberOfRetries = v } }

// OnError .
func OnError(v func(error)) func(o *options) { return func(o *options) { o.onError = v } }

// Period .
func Period(v time.Duration) func(o *options) { return func(o *options) { o.period = v } }

// Start starts an actor, close mailbox to signal stop
// and handler should respect that, or use another convention of ypurs.
func Start(
	handler func(<-chan payload) error,
	opt ...func(o *options)) (mailbox chan<- payload) {
	if handler == nil {
		panic("handler can not be nil")
	}
	o := new(options)
	for _, v := range opt {
		v(o)
	}
	if o.mailboxSize < 0 {
		o.mailboxSize = 0
	}
	q := make(chan payload, o.mailboxSize)
	go retry.Retry(
		func() error {
			return handler(q)
		},
		o.numberOfRetries,
		o.onError,
		o.period)
	mailbox = q
	return
}
