package mailbox

import "context"

// Storage like a []interface{} slice
type Storage interface {
	Len() int
	PeekHead() interface{}
	DropHead()
	Append(interface{})
}

// Mailbox an unbound mailbox that would not stop and close
// the receive channel until storage deplete
type Mailbox interface {
	Send() chan<- interface{}
	Receive() <-chan interface{}
}

type mailbox struct {
	ctx context.Context

	send    chan interface{}
	receive chan interface{}
	mails   Storage
}

func (mb *mailbox) loop() {
	var actualReceive chan interface{}
	first := func() interface{} {
		if mb.mails.Len() == 0 {
			return nil
		}
		return mb.mails.PeekHead()
	}
	for {
		select {
		case <-mb.ctx.Done():
			if mb.mails.Len() > 0 {
				continue
			}
			close(mb.receive)
			return
		case v := <-mb.send:
			mb.mails.Append(v)
			actualReceive = mb.receive
		case actualReceive <- first():
			mb.mails.DropHead()
			if mb.mails.Len() == 0 {
				actualReceive = nil
			}
		}
	}
}

func (mb *mailbox) Send() chan<- interface{}    { return mb.send }
func (mb *mailbox) Receive() <-chan interface{} { return mb.receive }

// New .
func New(ctx context.Context, store Storage) Mailbox {
	res := &mailbox{
		ctx:     ctx,
		send:    make(chan interface{}),
		receive: make(chan interface{}),
		mails:   store,
	}
	go res.loop()
	return res
}
